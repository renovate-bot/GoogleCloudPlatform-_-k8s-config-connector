// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runner

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	examples = `
# Run commands in multiple branches.
conductor runner --branch-repo=/usr/local/google/home/wfender/go/src/github.com/cheftako/k8s-config-connector_branches
`

	// flag names.
	branchConfigurationFlag = "branch-conf"
	mergeConfigurationFlag  = "merge-conf"
	branchRepoFlag          = "branch-repo"
	commandFlag             = "command"
	loggingDirFlag          = "logging-dir"
	readFileTypeFlag        = "file-type"
	controllerFilterFlag    = "controller-filter"
	testDirSuffixFlag       = "test-suffix"

	// Command values
	cmdDeleteGitBranch              = -5
	cmdHelp                         = 0
	cmdCheckRepo                    = 1
	cmdCreateGitBranch              = 2
	cmdMergeMetadata                = 3
	cmdEnableGCPAPIs                = 4
	cmdReadFiles                    = 5
	cmdWriteFiles                   = 6
	cmdDiff                         = 7
	cmdRevert                       = 8
	cmdPushBranch                   = 9 // New command for force pushing branches
	cmdCreateScriptYaml             = 10
	cmdCaptureHttpLog               = 11
	cmdGenerateMockGo               = 12
	cmdAddServiceRoundTrip          = 13
	cmdAddProtoMakefile             = 14
	cmdBuildProto                   = 15
	cmdCaptureMockOutput            = 16
	cmdRunAndFixMockTests           = 17
	cmdGenerateTypes                = 20
	cmdAdjustTypes                  = 21
	cmdGenerateCRD                  = 22
	cmdGenerateMapper               = 23
	cmdGenerateFuzzer               = 24
	cmdRunAndFixFuzzTests           = 25
	cmdRunAndFixAPIChecks           = 26
	cmdControllerClient             = 40
	cmdGenerateController           = 41
	cmdBuildAndFixController        = 42
	cmdCreateIdentity               = 43
	cmdControllerCreateTest         = 44
	cmdCaptureGoldenRealGCPOutput   = 45
	cmdRunAndFixGoldenRealGCPOutput = 46
	cmdCaptureGoldenMockOutput      = 47
	cmdRunAndFixGoldenMockOutput    = 48
	cmdMoveExistingTest             = 50
	cmdCreateFullTest               = 51

	typeScriptYaml = "scriptyaml"
	typeHttpLog    = "httplog"
	typeMockGo     = "mockgo"

	handleLocalChangeOptionCleanUp = "CLEANUP"
	handleLocalChangeOptionCommit  = "COMMIT"
	handleLocalChangeOptionFail    = "FAIL"
)

func BuildRunnerCmd() *cobra.Command {
	var opts RunnerOptions

	cmd := &cobra.Command{
		Use:     "runner",
		Short:   "runner Run commands in various branches.",
		Example: examples,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunRunner(cmd.Context(), &opts)
		},
		Args: cobra.ExactArgs(0),
	}

	cmd.Flags().StringVarP(&opts.branchConfFile, branchConfigurationFlag,
		"", "branches.yaml", "File containing the branch configurations.")
	cmd.Flags().StringVarP(&opts.branchMergeFile, mergeConfigurationFlag,
		"", "", "File containing the configurations to merge in.")
	cmd.Flags().StringVarP(&opts.branchRepoDir, branchRepoFlag,
		"", "", "Directory in which to do the work.")
	cmd.Flags().StringVarP(&opts.readFileType, readFileTypeFlag,
		"", "ScriptYaml", "Type of files to read, available values: 'ScriptYaml', 'HttpLog', 'MockGo'.")
	cmd.Flags().Int64VarP(&opts.command, commandFlag,
		"", 0, "Which commands to you on the directory.")
	cmd.Flags().StringVarP(&opts.loggingDir, loggingDirFlag,
		"", "", "dedicated directory for logging, empty for stdout.")
	cmd.Flags().DurationVarP(&opts.timeout, "timeout",
		"t", 5*time.Minute, "Global timeout for commands.")
	cmd.Flags().IntVarP(&opts.defaultRetries, "retries",
		"r", 10, "Default number of retries for failed commands.")
	cmd.Flags().StringVarP(&opts.forResources, "for-resources",
		"", "", "Comma-separated list of branch names to filter on.")
	cmd.Flags().StringVarP(&opts.forResourcesRegex, "for-resources-regex",
		"", "", "Regex to filter branch names to filter on.")
	cmd.Flags().BoolVarP(&opts.force, "force",
		"f", false, "Force operation even if files already exist.")
	cmd.Flags().IntVarP(&opts.numCommits, "num-commits",
		"n", 0, "Number of commits to diff/revert (default: 0)")
	cmd.Flags().BoolVarP(&opts.verbose, "verbose",
		"v", false, "Enable verbose output logging")
	cmd.Flags().StringVarP(&opts.branchSuffix, "branch-suffix",
		"", "", "Suffix to append to remote branch names when pushing")
	cmd.Flags().BoolVarP(&opts.skipMakeReadyPR, "skip-makereadypr",
		"", false, "Skip the make ready-pr step when pushing branches")
	cmd.Flags().StringVarP(&opts.processors, "processors",
		"", "", "Comma-separated list of processor function names to run (if empty, all processors are run)")
	cmd.Flags().StringVarP(&opts.controllerFilter, controllerFilterFlag,
		"", "", "Type of controller to filter for. (Eg terraform-v1beta1)")
	cmd.Flags().StringVarP(&opts.handleLocalChange, "handle-local-change",
		"", "", "Option to handle uncommitted local changes before switching to a different branch, available values: 'CLEANUP', 'COMMIT', 'FAIL'.")
	cmd.Flags().StringVarP(&opts.testDirSuffix, testDirSuffixFlag,
		"", "", "Suffix of the test to generate/run/fix for each branch")

	return cmd
}

type RunnerOptions struct {
	branchConfFile    string
	branchMergeFile   string // Secondary metadata file to merge into the primary
	branchRepoDir     string
	command           int64
	loggingDir        string
	timeout           time.Duration
	readFileType      string
	defaultRetries    int    // Default number of retries for commands
	forResources      string // Comma-separated list of branch names to filter on
	forResourcesRegex string // Regex to filter branches, only branches that match the regex are processed
	force             bool   // Force flag to override file existence checks
	numCommits        int    // Number of commits to diff (default: 1)
	verbose           bool   // Verbose output flag
	branchSuffix      string // Suffix to append to remote branch names when pushing
	skipMakeReadyPR   bool   // Skip make ready-pr step when pushing branches
	processors        string // Comma-separated list of processor function names to run
	controllerFilter  string // Filter the metadata for 1 type of controller. (Eg terraform-v1beta1)
	handleLocalChange string // Option to handle uncommitted local changes before switching to a different branch
	testDirSuffix     string // Suffix for test directory
}

func (opts *RunnerOptions) validateAndDefaultFlags() error {
	if opts.defaultRetries < 0 {
		return fmt.Errorf("retries flag cannot be negative, got %d", opts.defaultRetries)
	}

	switch opts.handleLocalChange {
	case "":
		// When the handle local change option is unset, each command will handle it differently.
		switch opts.command {
		case cmdCreateFullTest:
			opts.handleLocalChange = handleLocalChangeOptionCommit
		default:
			opts.handleLocalChange = handleLocalChangeOptionCleanUp
		}
	case handleLocalChangeOptionCleanUp, handleLocalChangeOptionCommit, handleLocalChangeOptionFail:
	default:
		return fmt.Errorf("handle-local-change flag must be set with one of %q, %q, %q but it is set to %q",
			handleLocalChangeOptionCleanUp, handleLocalChangeOptionCommit, handleLocalChangeOptionFail, opts.handleLocalChange)
	}

	if opts.testDirSuffix == "" {
		// When the test directory suffix is unset, each command will handle it differently.
		switch opts.command {
		case cmdControllerCreateTest, cmdCaptureGoldenRealGCPOutput,
			cmdRunAndFixGoldenRealGCPOutput, cmdCaptureGoldenMockOutput,
			cmdRunAndFixGoldenMockOutput:
			opts.testDirSuffix = "minimal"
		case cmdCreateFullTest:
			opts.testDirSuffix = "full"
		}
	} else {
		// Ensure the suffix is in lowercase.
		opts.testDirSuffix = strings.ToLower(opts.testDirSuffix)
	}
	return nil
}

type Branches struct {
	Branches []Branch `yaml:"branches"`
}

type Proto struct {
	Service   string `yaml:"service"`   // ai
	Package   string `yaml:"package"`   // google.ai.generativelanguage.v1beta
	Proto     string `yaml:"proto"`     // Model
	Kind      string `yaml:"kind"`      // AIModel
	ProtoPath string `yaml:"protopath"` // google.ai.generativelanguage.v1beta.Model
	Validated string `yaml:"validated"` // UNUSED
}

type Branch struct {
	Name       string `yaml:"name"`       // ai-model
	Local      string `yaml:"local"`      // resource-ai-model
	Remote     string `yaml:"remote"`     // resource-ai-model
	Command    string `yaml:"command"`    // gcloud ai models
	Group      string `yaml:"group"`      // ai
	Resource   string `yaml:"resource"`   // model
	Controller string `yaml:"controller"` // Unknown
	Skip       bool   `yaml:"skip"`       // Skip this branch in processing

	Kind      string `yaml:"kind"`          // AIModel
	Package   string `yaml:"package"`       // google.ai.generativelanguage.v1beta
	Proto     string `yaml:"proto"`         // Model
	ProtoPath string `yaml:"proto-path"`    // google.ai.generativelanguage.v1beta.model_service
	ProtoSvc  string `yaml:"proto-service"` // google.ai.generativelanguage.v1beta.ModelService
	ProtoMsg  string `yaml:"proto-msg"`     // google.ai.generativelanguage.v1beta.Model
	HostName  string `yaml:"host-name"`     // generativelanguage.googleapis.com

	Notes []string `yaml:"notes"` // Observation goes here

	// GCP APIs that need to be enabled for the branch.
	// example:
	// - eventarc.googleapis.com
	// - aiplatform.googleapis.com
	ApisEnabled []string `yaml:"apis-enabled"`
}

// BranchModifier is a function type that takes a Branch and returns a modified Branch
type BranchModifier func(opts *RunnerOptions, branch Branch, workDir string) Branch

func RunRunner(ctx context.Context, opts *RunnerOptions) error {
	log.Printf("Running conductor runner with branch config: %s", opts.branchConfFile)

	if opts.loggingDir == "" {
		return fmt.Errorf("logging-dir is required")
	}
	// If loggingDir is a relative path, make it absolute by prepending the current working directory
	if !filepath.IsAbs(opts.loggingDir) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current working directory: %w", err)
		}
		opts.loggingDir = filepath.Join(cwd, opts.loggingDir)
	}

	if err := opts.validateAndDefaultFlags(); err != nil {
		return err
	}

	log.Printf("Starting Runner")

	var branches Branches
	{
		b, err := os.ReadFile(opts.branchConfFile) // For read access.
		if err != nil {
			return fmt.Errorf("reading branch config file %q: %w", opts.branchConfFile, err)
		}

		if err := yaml.Unmarshal(b, &branches); err != nil {
			return fmt.Errorf("parsing branch config file %q: %w", opts.branchConfFile, err)
		}
	}

	// Filter based on forResources if specified
	if opts.forResources != "" {
		targetBranches := make(map[string]bool)
		for _, name := range strings.Split(opts.forResources, ",") {
			targetBranches[strings.TrimSpace(name)] = true
		}

		var resourceFilteredBranches []Branch
		for _, branch := range branches.Branches {
			if targetBranches[branch.Name] {
				resourceFilteredBranches = append(resourceFilteredBranches, branch)
			}
		}
		branches.Branches = resourceFilteredBranches
	}
	if opts.forResourcesRegex != "" {
		resourceFilteredBranches := make([]Branch, 0)
		for _, branch := range branches.Branches {
			match, err := regexp.MatchString(opts.forResourcesRegex, branch.Name)
			if err != nil {
				return fmt.Errorf("error matching regex %q: %w", opts.forResourcesRegex, err)
			}
			if match {
				resourceFilteredBranches = append(resourceFilteredBranches, branch)
			}
		}
		branches.Branches = resourceFilteredBranches
	}

	// Only filter skipped branches if not running metadata commands
	if opts.command >= 0 {
		// Filter out skipped branches
		var filteredBranches []Branch
		for _, branch := range branches.Branches {
			if !branch.Skip {
				filteredBranches = append(filteredBranches, branch)
			}
		}
		branches.Branches = filteredBranches

	}

	switch opts.command {
	case -6:
		if err := filterMetadata(opts, branches); err != nil {
			log.Fatalf("unnable to filter metadata: %v", err)
		}
	case cmdDeleteGitBranch: // -5
		for idx, branch := range branches.Branches {
			log.Printf("Delete GitHub Branch: %d name: %s, branch: %s\r\n", idx, branch.Name, branch.Local)
			deleteGithubBranch(opts, branch)
		}
	case -4:
		fixMetadata(opts, branches, "Skip(with reason)", setSkipOnBranchModifier)
	case -3:
		fixMetadata(opts, branches, "EnableAPIs", addEnableAPIsModifier)
	case -2:
		splitMetadata(opts, branches)
	case -1:
		fixMetadata(opts, branches, "ProtoPath", inferProtoPathModifier)
	case cmdHelp: // 0
		printHelp()
	case cmdCheckRepo: // 1
		checkRepoDir(ctx, opts, branches)
	case cmdCreateGitBranch: // 2
		for idx, branch := range branches.Branches {
			/*
				if branch.Controller != "Unknown" {
					// Skipping TF, DCL and Direct controller resources.
					continue
				}
			*/
			log.Printf("Create GitHub Branch: %d name: %s, branch: %s\r\n", idx, branch.Name, branch.Local)
			createGithubBranch(opts, branch)
		}
	case cmdMergeMetadata: // 3
		if err := mergeMetadata(opts, branches); err != nil {
			log.Fatalf("unnable to merge metadata: %v", err)
		}
	case cmdEnableGCPAPIs: // 4
		for idx, branch := range branches.Branches {
			log.Printf("Enable GCP APIs: %d name: %s, branch: %s\r\n", idx, branch.Name, branch.Local)
			enableAPIs(opts, branch)
		}
	case cmdReadFiles: // 5
		readFuncs := map[string]func(*RunnerOptions, Branch){
			typeScriptYaml: readScriptYaml,
			typeHttpLog:    readHttpLog,
			typeMockGo:     readMockGo,
		}

		if readFunc, ok := readFuncs[strings.ToLower(opts.readFileType)]; ok {
			for idx, branch := range branches.Branches {
				log.Printf("Read %s: %d name: %s, branch: %s\r\n", opts.readFileType, idx, branch.Name, branch.Local)
				readFunc(opts, branch)
			}
		}

	case cmdWriteFiles: // 6
		writeFuncs := map[string]func(*RunnerOptions, Branch){
			typeScriptYaml: writeScriptYaml,
		}

		if writeFunc, ok := writeFuncs[strings.ToLower(opts.readFileType)]; ok {
			for idx, branch := range branches.Branches {
				log.Printf("Write %s: %d name: %s, branch: %s\r\n", opts.readFileType, idx, branch.Name, branch.Local)
				writeFunc(opts, branch)
			}
		}

	case cmdDiff: // 7
		for idx, branch := range branches.Branches {
			log.Printf("Showing diff for last %d commits: %d name: %s, branch: %s\r\n", opts.numCommits, idx, branch.Name, branch.Local)
			// diffLastNCommits(opts, branch)
			determineCommandByCommit(opts, branch)
		}

	case cmdRevert: // 8
		if opts.numCommits <= 0 {
			log.Printf("Skipping revert, num-commits must be positive (got %d)", opts.numCommits)
			return nil
		}
		for idx, branch := range branches.Branches {
			log.Printf("Reverting last %d commits in branch %s: %d name: %s, branch: %s\r\n", opts.numCommits, branch.Name, idx, branch.Name, branch.Local)
			revertLastNCommits(opts, branch)
		}
	case cmdPushBranch: // 9
		processors := []BranchProcessor{
			{Fn: runAPIChecks, CommitMsgTemplate: COMMIT_MSG_9A, AttemptsOnNoChange: 1},
			{Fn: makeReadyPR, CommitMsgTemplate: COMMIT_MSG_9B, AttemptsOnChanges: 5, AttemptsOnNoChange: 1},
			{Fn: pushBranch, CommitMsgTemplate: COMMIT_MSG_9C, AttemptsOnNoChange: 1},
		}
		processBranches(ctx, opts, REGEX_MSG_9C, branches.Branches, "Prep and Push Branch", processors)
	case cmdCreateScriptYaml: // 10
		processBranches(ctx, opts, REGEX_MSG_10, branches.Branches, "Script YAML", []BranchProcessor{{Fn: createScriptYaml, CommitMsgTemplate: COMMIT_MSG_10}})
	case cmdCaptureHttpLog: // 11
		processBranches(ctx, opts, REGEX_MSG_11, branches.Branches, "HTTP Log", []BranchProcessor{{Fn: captureHttpLog, CommitMsgTemplate: COMMIT_MSG_11}})
	case cmdGenerateMockGo: // 12
		processBranches(ctx, opts, REGEX_MSG_12, branches.Branches, "Mock Go Files", []BranchProcessor{{Fn: generateMockGo, CommitMsgTemplate: COMMIT_MSG_12}})
	case cmdAddServiceRoundTrip: // 13
		processBranches(ctx, opts, REGEX_MSG_13, branches.Branches, "Service RoundTrip", []BranchProcessor{{Fn: addServiceToRoundTrip, CommitMsgTemplate: COMMIT_MSG_13}})
	case cmdAddProtoMakefile: // 14
		processBranches(ctx, opts, REGEX_MSG_14, branches.Branches, "Proto Makefile", []BranchProcessor{{Fn: addProtoToMakefile, CommitMsgTemplate: COMMIT_MSG_14, AttemptsOnNoChange: 1}})
	case cmdBuildProto: // 15
		processBranches(ctx, opts, REGEX_MSG_15, branches.Branches, "Build Proto", []BranchProcessor{{Fn: buildProtoFiles, CommitMsgTemplate: COMMIT_MSG_15}})
	case cmdCaptureMockOutput: // 16
		processBranches(ctx, opts, REGEX_MSG_16, branches.Branches, "Mock Tests", []BranchProcessor{{Fn: runMockgcpTests, CommitMsgTemplate: COMMIT_MSG_16}})
	case cmdRunAndFixMockTests: // 17
		processBranches(ctx, opts, REGEX_MSG_17, branches.Branches, "Mock Tests", []BranchProcessor{{Fn: fixMockgcpFailures, CommitMsgTemplate: COMMIT_MSG_17, VerifyFn: runMockgcpTests, VerifyAttempts: 10, AttemptsOnNoChange: 2}})
	case cmdGenerateTypes: // 20
		processBranches(ctx, opts, REGEX_MSG_20, branches.Branches, "Types", []BranchProcessor{{Fn: generateTypes, CommitMsgTemplate: COMMIT_MSG_20}})
	case cmdAdjustTypes: // 21
		processors := []BranchProcessor{
			{Fn: setTypeSpecStatus, CommitMsgTemplate: COMMIT_MSG_21A, AttemptsOnNoChange: 6, SkipProcessorOnMsg: skipPost21A},
			{Fn: setTypeParent, CommitMsgTemplate: COMMIT_MSG_21B, AttemptsOnNoChange: 6, SkipProcessorOnMsg: skipPost21B},
			{Fn: adjustIdentityParent, CommitMsgTemplate: COMMIT_MSG_21C, AttemptsOnNoChange: 6, SkipProcessorOnMsg: skipPost21C},
			{Fn: regenerateTypes, CommitMsgTemplate: COMMIT_MSG_21D, SkipProcessorOnMsg: skipPost21D},
			{Fn: removeNameField, CommitMsgTemplate: COMMIT_MSG_21E, SkipProcessorOnMsg: skipPost21E},
			{Fn: moveEtagField, CommitMsgTemplate: COMMIT_MSG_21F, AttemptsOnNoChange: 1, SkipProcessorOnMsg: skipPost21F},
			// add a kubebuilder required field label to the fields that are marked required in proto
			{Fn: addRequiredFieldTags, CommitMsgTemplate: COMMIT_MSG_21G},
			// TODO? manual: Add something to handle references to other resources: https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/4010/commits/1651a0a7af5bca37b5c2e134dd3f600ebac6a172
			// * https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/4017/commits/cc726106aff55d41e6bc94272acc3612f2636397
		}
		processBranches(ctx, opts, REGEX_MSG_21G, branches.Branches, "Adjusting types", processors)
	case cmdGenerateCRD: // 22
		processBranches(ctx, opts, REGEX_MSG_22, branches.Branches, "CRD", []BranchProcessor{{Fn: generateCRD, CommitMsgTemplate: COMMIT_MSG_22}})
	case cmdGenerateMapper: // 23
		processBranches(ctx, opts, REGEX_MSG_23, branches.Branches, "Mapper", []BranchProcessor{{Fn: generateMapper, CommitMsgTemplate: COMMIT_MSG_23}})
	case cmdGenerateFuzzer: // 24
		processBranches(ctx, opts, REGEX_MSG_24, branches.Branches, "Fuzzer", []BranchProcessor{{Fn: generateFuzzer, CommitMsgTemplate: COMMIT_MSG_24}})
	case cmdRunAndFixFuzzTests: // 25
		processBranches(ctx, opts, REGEX_MSG_25, branches.Branches, "Fuzzer Tests", []BranchProcessor{{Fn: fixFuzzerFailures, CommitMsgTemplate: COMMIT_MSG_25, VerifyFn: runFuzzerTests, VerifyAttempts: 5, AttemptsOnNoChange: 2}})
	case cmdRunAndFixAPIChecks: // 26
		processBranches(ctx, opts, REGEX_MSG_26, branches.Branches, "API Checks", []BranchProcessor{{Fn: fixAPICheckFailures, CommitMsgTemplate: COMMIT_MSG_26, VerifyFn: runAPIChecks, VerifyAttempts: 5}})
	case cmdControllerClient: // 40
		processBranches(ctx, opts, REGEX_MSG_40, branches.Branches, "Controller Client", []BranchProcessor{{Fn: generateControllerClient, CommitMsgTemplate: COMMIT_MSG_40}})
	case cmdGenerateController: // 41
		processBranches(ctx, opts, REGEX_MSG_41, branches.Branches, "Controller", []BranchProcessor{{Fn: generateController, CommitMsgTemplate: COMMIT_MSG_41}})
	case cmdBuildAndFixController: // 42
		processBranches(ctx, opts, REGEX_MSG_42, branches.Branches, "Build and Fix Controller", []BranchProcessor{{Fn: fixControllerBuild, CommitMsgTemplate: COMMIT_MSG_42, VerifyFn: buildController, VerifyAttempts: 10}})
	case cmdCreateIdentity: // 43
		processBranches(ctx, opts, REGEX_MSG_43B, branches.Branches, "Identity and Reference", []BranchProcessor{
			{Fn: generateControllerIdentity, CommitMsgTemplate: COMMIT_MSG_43A},
			{Fn: generateControllerReference, CommitMsgTemplate: COMMIT_MSG_43B},
		})
	case cmdControllerCreateTest: // 44
		processBranches(ctx, opts, REGEX_MSG_44B, branches.Branches, "Controller Test", []BranchProcessor{
			{Fn: createControllerTest, CommitMsgTemplate: COMMIT_MSG_44A},
			{Fn: updateTestHarness, CommitMsgTemplate: COMMIT_MSG_44B},
		})
	case cmdCaptureGoldenRealGCPOutput: // 45
		processBranches(ctx, opts, REGEX_MSG_45, branches.Branches, "Record Golden Real GCP Tests", []BranchProcessor{{Fn: runGoldenRealGCPTests, CommitMsgTemplate: COMMIT_MSG_45, AttemptsOnNoChange: 1}})
	case cmdRunAndFixGoldenRealGCPOutput: // 46
		processBranches(ctx, opts, REGEX_MSG_46, branches.Branches, "Fix Real GCP Tests", []BranchProcessor{{Fn: fixGoldenTests, CommitMsgTemplate: COMMIT_MSG_46, VerifyFn: runGoldenRealGCPTests, VerifyAttempts: 5, AttemptsOnNoChange: 2}})
	case cmdCaptureGoldenMockOutput: // 47
		processBranches(ctx, opts, REGEX_MSG_47, branches.Branches, "Record Golden Mock Tests", []BranchProcessor{{Fn: runGoldenMockTests, CommitMsgTemplate: COMMIT_MSG_47, AttemptsOnNoChange: 1}})
	case cmdRunAndFixGoldenMockOutput: // 48
		processBranches(ctx, opts, REGEX_MSG_48, branches.Branches, "Fix Mock GCP Tests", []BranchProcessor{{Fn: fixMockGcpForGoldenTests, CommitMsgTemplate: COMMIT_MSG_48, VerifyFn: runGoldenMockTests, VerifyAttempts: 5, AttemptsOnNoChange: 2}})
	case cmdMoveExistingTest: // 50
		processBranches(ctx, opts, REGEX_MSG_50, branches.Branches, "Move Existing Test", []BranchProcessor{{Fn: moveTestToSubDir, CommitMsgTemplate: COMMIT_MSG_50, AttemptsOnNoChange: 0, CommitOptional: true}})
	case cmdCreateFullTest: // 51
		processBranches(ctx, opts, REGEX_MSG_51, branches.Branches, "Create Maximal Test", []BranchProcessor{
			{Fn: moveTestToSubDir, CommitMsgTemplate: COMMIT_MSG_51A, AttemptsOnNoChange: 0, CommitOptional: true, SkipProcessorOnMsg: skipPost51A},
			{Fn: createFullCoverageTest, CommitMsgTemplate: COMMIT_MSG_51B, CommitOptional: true},
			{Fn: updateTestHarness, CommitMsgTemplate: COMMIT_MSG_51C, CommitOptional: true},
		})
	default:
		log.Fatalf("unrecognized command: %d", opts.command)
	}
	return nil
}

func printHelp() {
	log.Println("conductor runner --branch-repo=? --branch-conf=< META > --command=<CMD>")
	log.Println("\t<CMD>")
	log.Println("\t0 - Print help")
	log.Println("\t1 - [Validate] Repo directory and metadata")
	log.Println("\t2 - [Branch] Create the local github branches from the metadata")
	log.Println("\t3 - [Branch] Merge updated metadata file into main metadata file")
	log.Println("\t4 - [Project] Enable GCP APIs for each branch")
	log.Println("\t5 - [Generated] Read the specific type of generated files in each github branch")
	log.Println("\t6 - [Generated] Write the specific type of files from all_scripts.yaml to each github branch")
	log.Println("\t7 - [Git] Use git to determine the last command run on each branch")
	log.Println("\t8 - [Git] Revert last N commits in each branch (use -n to specify N)")
	log.Println("\t9 - [Git] Make ready-pr and push each branch to origin. Use --branch-suffix to specify a suffix for the remote branch if needed")
	log.Println("\t10 - [Mock] Create script.yaml for mock gcp generation in each github branch")
	log.Println("\t11 - [Mock] Create _http.log for mock gcp generation in each github branch")
	log.Println("\t12 - [Mock] Generate mock Service and Resource go files in each github branch")
	log.Println("\t13 - [Mock] Add service to mock_http_roundtrip.go in each github branch")
	log.Println("\t14 - [Mock] Add proto to makefile in each github branch")
	log.Println("\t15 - [Proto] Build proto files in mockgcp directory")
	log.Println("\t16 - [Mock] Capture mock golden output for each branch")
	log.Println("\t17 - [Mock] Run and Fix mockgcp tests in each github branch")
	log.Println("\t20 - [CRD] Generate Types for each branch")
	log.Println("\t21 - [CRD] Adjust the types for each branch")
	log.Println("\t22 - [CRD] Generate CRD for each branch")
	log.Println("\t23 - [CRD] Generate Mapper for each branch")
	log.Println("\t24 - [CRD] Generate Fuzzer for each branch")
	log.Println("\t25 - [CRD] Run and Fix fuzzer tests for each branch")
	log.Println("\t26 - [CRD] Run and Fix API checks for each branch")
	log.Println("\t40 - [Controller] Generate controller client for each branch")
	log.Println("\t41 - [Controller] Generate controller for each branch")
	log.Println("\t42 - [Controller] Build and fix controller for each branch")
	log.Println("\t43 - [Controller] [optional, simialr to 20, 21] Create identity and reference files for each branch")
	log.Println("\t44 - [Controller] Create minimal test files for each branch")
	log.Println("\t45 - [Controller] Capture golden logs for real GCP tests for each branch")
	log.Println("\t46 - [Controller] Run and Fix real GCP tests for each branch")
	log.Println("\t47 - [Controller] Capture golden logs for mock GCP tests for each branch")
	log.Println("\t48 - [Controller] Run and Fix mock GCP tests for each branch")
	log.Println("\t50 - [Test] Move test data to subdirectory if the test files are directly under <version>/<kind> directory for each branch")
	log.Println("\t51 - [Test] Create a maximal test for each branch")
}

func checkRepoDir(ctx context.Context, opts *RunnerOptions, branches Branches) {
	cmd := exec.CommandContext(ctx, "ls", "-alh")
	cmd.Dir = opts.branchRepoDir

	results, err := execCommand(cmd)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ls-alh: %s", results.Stdout)

	// Check for uniqueness constraints in the metadata.
	gcloudMap := make(map[string]string)
	nameMap := make(map[string]Branch)
	gitMap := make(map[string]string)
	grMap := make(map[string]Branch)
	kindMap := make(map[string]string)
	for idx, branch := range branches.Branches {
		if branch.Command != "" {
			if existing, ok := gcloudMap[branch.Command]; ok {
				log.Printf("Command (%s) uniqueness constraint between %s and (%d)%s\r",
					branch.Command, existing, idx, branch.Name)
			}
			gcloudMap[branch.Command] = branch.Name
		}
		if existing, ok := nameMap[branch.Name]; ok {
			log.Printf("Name uniqueness constraint between %s at and %s\r",
				branch.Name, existing.Name)
		}
		nameMap[branch.Name] = branch

		if existing, ok := gitMap[branch.Local]; ok {
			log.Printf("Github uniqueness constraint between %s and (%d)%s\r",
				existing, idx, branch.Name)
		}
		gitMap[branch.Local] = branch.Name

		gr := branch.Group + ":" + branch.Resource
		if existing, ok := grMap[gr]; ok {
			log.Printf("Branch:Resource uniqueness constraint between %s and (%d)%s\r",
				existing.Name, idx, branch.Name)
		}
		grMap[gr] = branch

		if branch.Kind != "" {
			if existing, ok := kindMap[branch.Kind]; ok {
				log.Printf("Kind uniqueness (%s) constraint between names: %s and (%d)%s\r",
					branch.Kind, existing, idx, branch.Name)
			}
		}
		kindMap[branch.Kind] = branch.Name
	}

	// Fix the data and write back
	/*
		var newBranches Branches
		// newBranches.Branches = make([]Branch, len(branches.Branches))
		for _, branch := range branches.Branches {
			branch.Name = strings.ToLower(branch.Name)
			branch.Local = strings.ToLower(branch.Local)
			branch.Remote = strings.ToLower(branch.Remote)
			branch.Group = strings.ToLower(branch.Group)
			branch.Resource = strings.ToLower(branch.Resource)
			branch.Notes = make([]string, 1)
			branch.Notes[0] = strings.TrimSpace(branch.Note)
			branch.Note = ""
			branch.Controller = "Unknown"
			newBranches.Branches = append(newBranches.Branches, branch)
		}
		data, err := yaml.Marshal(newBranches)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile("branches-new.yaml", data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	*/

	/*
		// Fold in the PROTO file
		file, err := os.Open("proto-list.yaml") // For read access.
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Read file %s: context type %T", "proto-list.yaml", string(data))
		rawProtos := strings.Split(string(data), "---")
		var proto Proto
		var newBranches Branches
		var filledNames []string
		for _, rawProto := range rawProtos {
			err = yaml.Unmarshal([]byte(rawProto), &proto)
			if err != nil {
				log.Printf("Failed to parse %v\r", rawProto)
				log.Fatalf("error: %v", err)
			}
			kindGr := proto.Service + ":" + strings.ToLower(proto.Kind)
			protoGr := proto.Service + ":" + strings.ToLower(proto.Proto)
			kindName := proto.Service + "-" + strings.ToLower(proto.Kind)
			protoName := proto.Service + "-" + strings.ToLower(proto.Proto)
			if branch, ok := grMap[kindGr]; ok {
				branch.Kind = proto.Kind
				branch.Package = proto.Package
				branch.Proto = proto.Proto
				branch.ProtoPath = proto.ProtoPath
				newBranches.Branches = append(newBranches.Branches, branch)
				filledNames = append(filledNames, branch.Name)
			} else if branch, ok := grMap[protoGr]; ok {
				branch.Kind = proto.Kind
				branch.Package = proto.Package
				branch.Proto = proto.Proto
				branch.ProtoPath = proto.ProtoPath
				newBranches.Branches = append(newBranches.Branches, branch)
				filledNames = append(filledNames, branch.Name)
			} else if branch, ok := nameMap[kindName]; ok {
				branch.Kind = proto.Kind
				branch.Package = proto.Package
				branch.Proto = proto.Proto
				branch.ProtoPath = proto.ProtoPath
				newBranches.Branches = append(newBranches.Branches, branch)
				filledNames = append(filledNames, branch.Name)
			} else if branch, ok := nameMap[protoName]; ok {
				branch.Kind = proto.Kind
				branch.Package = proto.Package
				branch.Proto = proto.Proto
				branch.ProtoPath = proto.ProtoPath
				newBranches.Branches = append(newBranches.Branches, branch)
				filledNames = append(filledNames, branch.Name)
			} else {
				log.Printf("No match found for %s or %s\r", kindName, protoName)
				branch.Name = protoName
				branch.Local = "resource-" + protoName
				branch.Remote = "resource-" + protoName
				branch.Group = proto.Service
				branch.Resource = proto.Kind
				branch.Controller = "Unknown"
				branch.Notes = make([]string, 1)
				branch.Notes[0] = "No gcloud command found"
				branch.Kind = proto.Kind
				branch.Package = proto.Package
				branch.Proto = proto.Proto
				branch.ProtoPath = proto.ProtoPath
				newBranches.Branches = append(newBranches.Branches, branch)
				filledNames = append(filledNames, branch.Name)
			}
		}
		for _, branch := range branches.Branches {
			if slices.Contains(filledNames, branch.Name) {
				continue
			}
			newBranches.Branches = append(newBranches.Branches, branch)
			log.Printf("No additional metadata found for %s\r", branch.Name)
		}
		log.Printf("Marshalling %d branches\r", len(branches.Branches))
		data, err = yaml.Marshal(newBranches)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile("branches-new.yaml", data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	*/
}

func fixMetadata(opts *RunnerOptions, branches Branches, what string, modifier BranchModifier) {
	workDir := filepath.Join(opts.branchRepoDir, ".build", "third_party", "googleapis")
	var newBranches Branches
	for _, branch := range branches.Branches {
		log.Printf("Fixing Metadata: %s for branch: %s", what, branch.Name)
		newBranch := modifier(opts, branch, workDir)
		newBranches.Branches = append(newBranches.Branches, newBranch)
	}
	writeBranchesStableOrder(newBranches, "branches-new.yaml")
}

func inferProtoPathModifier(opts *RunnerOptions, branch Branch, workDir string) Branch {
	if branch.ProtoPath == "" {
		branch.ProtoPath = inferProtoPath(opts, branch, workDir)
		log.Printf("ProtoPath for %s should be %s", branch.Name, branch.ProtoPath)
	}
	return branch
}

func addEnableAPIsModifier(opts *RunnerOptions, branch Branch, workDir string) Branch {
	close := setLoggingWriter(opts, branch)
	defer close()
	if branch.Command == "" {
		return branch
	}
	if branch.ApisEnabled != nil && len(branch.ApisEnabled) > 0 {
		return branch
	}

	cfg := CommandConfig{
		Name:    "API Discovery",
		Cmd:     "codebot",
		Args:    []string{"--prompt=/dev/stdin"},
		WorkDir: workDir,
		Stdin: strings.NewReader(fmt.Sprintf(`Given the gcloud command %q, what Google Cloud APIs need to be enabled to use this command?
Try inferring the apis required from the command.
If needed use gcloud command to instrospect the error message to get the API name.
Please respond with a list of API service names in the format needed for 'gcloud services enable'.
For example, if a command needs the Cloud Storage API, respond with 'storage.googleapis.com'.
If you are using gcloud to create something, make sure you are deleting it after you are done.
Print the list of APIs, one on each line in this format api-required: <api-name>
Only include APIs that are directly needed by this command.
`, branch.Command)),
		RetryBackoff: GenerativeCommandRetryBackoff,
	}

	output, err := executeCommand(opts, cfg)
	if err != nil {
		log.Printf("Failed to get APIs for %s: %v", branch.Name, err)
		return branch
	}

	// Parse the codebot output to get API list
	var filteredLines []string
	for _, line := range strings.Split(strings.TrimSpace(output.Stdout), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "api-required:") {
			filteredLines = append(filteredLines, line)
		}
	}

	// split the structred output into a list of APIs
	var cleanedAPIs []string
	for _, api := range filteredLines {
		api = strings.TrimSpace(api)
		api = strings.TrimPrefix(api, "api-required:")
		api = strings.TrimSpace(api)
		if api != "" {
			cleanedAPIs = append(cleanedAPIs, api)
		}
	}

	branch.ApisEnabled = cleanedAPIs
	return branch
}

func inferProtoPath(opts *RunnerOptions, branch Branch, workDir string) string {
	var protoDir = ""
	var svcNm = ""
	if branch.Proto != "" {
		svcNm = branch.Proto
	} else if branch.ProtoSvc != "" {
		dirs := strings.Split(branch.ProtoSvc, ".")
		svcNm = dirs[len(dirs)-1]
	} else if branch.Resource != "" {
		svcNm = branch.Resource
	} else {
		return ""
	}

	if branch.ProtoPath != "" {
		dirs := strings.Split(branch.ProtoPath, ".")
		protoDir = filepath.Join(dirs[:len(dirs)-1]...)
	} else if branch.ProtoSvc != "" {
		dirs := strings.Split(branch.ProtoSvc, ".")
		protoDir = filepath.Join(dirs[:len(dirs)-1]...)
	} else if branch.Package != "" {
		dirs := strings.Split(branch.Package, ".")
		protoDir = filepath.Join(dirs...)
	} else {
		return ""
	}

	searchPath := filepath.Join(workDir, protoDir, "*.proto")
	files, err := filepath.Glob(searchPath)
	if err != nil {
		log.Printf("Glob error %v", err)
		return ""
	}
	if len(files) == 1 {
		localFile, _ := strings.CutPrefix(files[0], workDir+string(filepath.Separator))
		log.Printf("Glob for %s matched %s", branch.Name, localFile)
		return localFile
	}

	searchList := ""
	first := true
	args := []string{"-iH", fmt.Sprintf("^service %s", svcNm)}
	for _, file := range files {
		localFile, _ := strings.CutPrefix(file, workDir+string(filepath.Separator))
		args = append(args, localFile)
		if first {
			searchList = localFile
			first = false
		} else {
			searchList += " " + localFile
		}
	}

	cfg := CommandConfig{
		Name:        "Search for service",
		Cmd:         "egrep",
		Args:        args,
		WorkDir:     workDir,
		MaxAttempts: 2,
	}
	output, err := executeCommand(opts, cfg)
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			args := []string{"-iH", "^service"}
			for _, file := range files {
				localFile, _ := strings.CutPrefix(file, workDir+string(filepath.Separator))
				args = append(args, localFile)
				if first {
					searchList = localFile
					first = false
				} else {
					searchList += " " + localFile
				}
			}

			cfg = CommandConfig{
				Name:        "Search for any service",
				Cmd:         "egrep",
				Args:        args,
				WorkDir:     workDir,
				MaxAttempts: 2,
			}
			output, err = executeCommand(opts, cfg)
			if err != nil {
				log.Printf("Working in directory %s", workDir)
				log.Printf("Got response2 %v", output.Stderr)
				log.Printf("Find proto file error: %q\n", err)
				return ""
			}
		} else {
			log.Printf("Working in directory %s", workDir)
			log.Printf("Find proto file error: %q\n", err)
			return ""
		}
	}

	vals := strings.Split(output.Stdout, ":")
	if len(vals) <= 1 {
		log.Printf("ERROR: something wrong with grep response: %q\n", output.Stdout)
		return ""
	}
	return vals[0]
}

// Trying to put all resources with the same group together while
// keeping the buckets roughly the same size.
func splitMetadata(opts *RunnerOptions, branches Branches) {
	var newBranches [7]Branches
	groupSplitMap := make(map[string]int)
	for _, branch := range branches.Branches {
		if branch.Command == "" {
			continue
		}
		if split, present := groupSplitMap[branch.Group]; present {
			newBranches[split].Branches = append(newBranches[split].Branches, branch)
			continue
		}
		smallest := len(newBranches[0].Branches)
		bucket := 0
		for cntr := 1; cntr < 7; cntr++ {
			if len(newBranches[cntr].Branches) < smallest {
				smallest = len(newBranches[cntr].Branches)
				bucket = cntr
			}
		}
		newBranches[bucket].Branches = append(newBranches[bucket].Branches, branch)
		groupSplitMap[branch.Group] = bucket
	}
	// Hard coding splitting into 7 files
	for cntr := 0; cntr < 7; cntr++ {
		writeBranchesStableOrder(newBranches[cntr], fmt.Sprintf("branches-%d.yaml", cntr))
	}
}

// Filter metadata for a particular controller type.
func filterMetadata(opts *RunnerOptions, branches Branches) error {
	if opts.controllerFilter == "" {
		return fmt.Errorf("filterMetadata requires that controller fitler be set")
	}

	var filteredBranch Branches
	for _, branch := range branches.Branches {
		if branch.Controller != opts.controllerFilter {
			continue
		}
		filteredBranch.Branches = append(filteredBranch.Branches, branch)
	}
	writeBranchesStableOrder(filteredBranch, fmt.Sprintf("branches-%s.yaml", opts.controllerFilter))
	return nil
}

func mergeMetadata(opts *RunnerOptions, branches Branches) error {
	// Load the merge file
	var newBranches Branches
	{
		b, err := os.ReadFile(opts.branchMergeFile) // For read access.
		if err != nil {
			return fmt.Errorf("reading branch config file %q: %w", opts.branchMergeFile, err)
		}

		if err := yaml.Unmarshal(b, &newBranches); err != nil {
			return fmt.Errorf("parsing branch config file %q: %w", opts.branchMergeFile, err)
		}
	}
	log.Printf("Existing branches is size %d, merging branches is size %d", len(branches.Branches), len(newBranches.Branches))

	// Ensure skip is not set in the merged data
	// Also build the lookup table
	branchPresent := make(map[string]bool)
	for idx, branch := range newBranches.Branches {
		if branch.Skip == true {
			branch.Skip = false
			newBranches.Branches[idx] = branch
		}
		branchPresent[branch.Name] = true
	}

	// Merge the new data into the old.
	// We do this by copying in old data when absent
	for _, branch := range branches.Branches {
		if branchPresent[branch.Name] {
			continue
		}
		newBranches.Branches = append(newBranches.Branches, branch)
	}

	// Write out the finished data
	writeBranchesStableOrder(newBranches, "branches-new.yaml")
	return nil
}

func setSkipOnBranchModifier(opts *RunnerOptions, branch Branch, workDir string) Branch {
	close := setLoggingWriter(opts, branch)
	defer close()

	if branch.Command == "" {
		branch.Skip = true
		branch.Notes = append(branch.Notes, "no gcloud command")
	}

	// Run gcloud xxx --help to see if it implements CRUD methods
	gcloudCommand := branch.Command
	gcloudCommand = strings.TrimPrefix(gcloudCommand, "gcloud")
	args := strings.Fields(gcloudCommand)

	// Skip if gcloud command has no args
	if len(args) == 0 {
		branch.Skip = true
		branch.Notes = append(branch.Notes, "invalid gcloud command")
	}

	// Keep if any of CRUD is supported
	// args = append(args, "--help")
	cfg := CommandConfig{
		Name:        "Check CRUD Support",
		Cmd:         "gcloud",
		Args:        args,
		WorkDir:     workDir,
		MaxAttempts: 1,
	}
	op, _ := executeCommand(opts, cfg)
	// todo: shall we keep when gcloud xxx --help returns error?
	// Just keep everything for now
	// todo: Do we need to include 'patch'? I thought 'gcloud' primarily uses 'update' for these operations.
	var operationRegex = regexp.MustCompile(`(create|update|delete)`)
	if !operationRegex.MatchString(op.Stderr) {
		branch.Skip = true
		branch.Notes = append(branch.Notes, "gcloud command does not contain any of create/update/delete")
	}
	return branch
}

func determineCommandByCommit(opts *RunnerOptions, branch Branch) {
	close := setLoggingWriter(opts, branch)
	defer close()
	workDir := opts.branchRepoDir

	ctx := context.TODO()
	checkoutBranch(ctx, branch, workDir)

	// Run git diff command
	message, found := getLatestGitMessage(workDir, opts, branch)
	if !found {
		return
	}

	switch {
	case REGEX_MSG_9A.MatchString(message):
		log.Printf("%s: command 9A\n", branch.Name)
		fmt.Printf("%s: command 9A\n", branch.Name)
	case REGEX_MSG_9B.MatchString(message):
		log.Printf("%s: command 9B\n", branch.Name)
		fmt.Printf("%s: command 9B\n", branch.Name)
	case REGEX_MSG_9C.MatchString(message):
		log.Printf("%s: command 9C\n", branch.Name)
		fmt.Printf("%s: command 9C\n", branch.Name)
	case REGEX_MSG_10.MatchString(message):
		log.Printf("%s: command 10\n", branch.Name)
		fmt.Printf("%s: command 10\n", branch.Name)
	case REGEX_MSG_11.MatchString(message):
		log.Printf("%s: command 11\n", branch.Name)
		fmt.Printf("%s: command 11\n", branch.Name)
	case REGEX_MSG_12.MatchString(message):
		log.Printf("%s: command 12\n", branch.Name)
		fmt.Printf("%s: command 12\n", branch.Name)
	case REGEX_MSG_13.MatchString(message):
		log.Printf("%s: command 13\n", branch.Name)
		fmt.Printf("%s: command 13\n", branch.Name)
	case REGEX_MSG_14.MatchString(message):
		log.Printf("%s: command 14\n", branch.Name)
		fmt.Printf("%s: command 14\n", branch.Name)
	case REGEX_MSG_15.MatchString(message):
		log.Printf("%s: command 15\n", branch.Name)
		fmt.Printf("%s: command 15\n", branch.Name)
	case REGEX_MSG_16.MatchString(message):
		log.Printf("%s: command 16\n", branch.Name)
		fmt.Printf("%s: command 16\n", branch.Name)
	case REGEX_MSG_17.MatchString(message):
		log.Printf("%s: command 17\n", branch.Name)
		fmt.Printf("%s: command 17\n", branch.Name)
	case REGEX_MSG_20.MatchString(message):
		log.Printf("%s: command 20\n", branch.Name)
		fmt.Printf("%s: command 20\n", branch.Name)
	case REGEX_MSG_21A.MatchString(message):
		log.Printf("%s: command 21A\n", branch.Name)
		fmt.Printf("%s: command 21A\n", branch.Name)
	case REGEX_MSG_21B.MatchString(message):
		log.Printf("%s: command 21B\n", branch.Name)
		fmt.Printf("%s: command 21B\n", branch.Name)
	case REGEX_MSG_21C.MatchString(message):
		log.Printf("%s: command 21C\n", branch.Name)
		fmt.Printf("%s: command 21C\n", branch.Name)
	case REGEX_MSG_21D.MatchString(message):
		log.Printf("%s: command 21D\n", branch.Name)
		fmt.Printf("%s: command 21D\n", branch.Name)
	case REGEX_MSG_21E.MatchString(message):
		log.Printf("%s: command 21E\n", branch.Name)
		fmt.Printf("%s: command 21E\n", branch.Name)
	case REGEX_MSG_21F.MatchString(message):
		log.Printf("%s: command 21F\n", branch.Name)
		fmt.Printf("%s: command 21F\n", branch.Name)
	case REGEX_MSG_21G.MatchString(message):
		log.Printf("%s: command 21G\n", branch.Name)
		fmt.Printf("%s: command 21G\n", branch.Name)
	case REGEX_MSG_22.MatchString(message):
		log.Printf("%s: command 22\n", branch.Name)
		fmt.Printf("%s: command 22\n", branch.Name)
	case REGEX_MSG_23.MatchString(message):
		log.Printf("%s: command 23\n", branch.Name)
		fmt.Printf("%s: command 23\n", branch.Name)
	case REGEX_MSG_24.MatchString(message):
		log.Printf("%s: command 24\n", branch.Name)
		fmt.Printf("%s: command 24\n", branch.Name)
	case REGEX_MSG_25.MatchString(message):
		log.Printf("%s: command 25\n", branch.Name)
		fmt.Printf("%s: command 25\n", branch.Name)
	case REGEX_MSG_26.MatchString(message):
		log.Printf("%s: command 26\n", branch.Name)
		fmt.Printf("%s: command 26\n", branch.Name)
	case REGEX_MSG_40.MatchString(message):
		log.Printf("%s: command 40\n", branch.Name)
		fmt.Printf("%s: command 40\n", branch.Name)
	case REGEX_MSG_41.MatchString(message):
		log.Printf("%s: command 41\n", branch.Name)
		fmt.Printf("%s: command 41\n", branch.Name)
	case REGEX_MSG_42.MatchString(message):
		log.Printf("%s: command 42\n", branch.Name)
		fmt.Printf("%s: command 42\n", branch.Name)
	case REGEX_MSG_43A.MatchString(message):
		log.Printf("%s: command 43A\n", branch.Name)
		fmt.Printf("%s: command 43A\n", branch.Name)
	case REGEX_MSG_43B.MatchString(message):
		log.Printf("%s: command 43B\n", branch.Name)
		fmt.Printf("%s: command 43B\n", branch.Name)
	case REGEX_MSG_44A.MatchString(message):
		log.Printf("%s: command 44A\n", branch.Name)
		fmt.Printf("%s: command 44A\n", branch.Name)
	case REGEX_MSG_44B.MatchString(message):
		log.Printf("%s: command 44B\n", branch.Name)
		fmt.Printf("%s: command 44B\n", branch.Name)
	case REGEX_MSG_45.MatchString(message):
		log.Printf("%s: command 45\n", branch.Name)
		fmt.Printf("%s: command 45\n", branch.Name)
	case REGEX_MSG_46.MatchString(message):
		log.Printf("%s: command 46\n", branch.Name)
		fmt.Printf("%s: command 46\n", branch.Name)
	case REGEX_MSG_47.MatchString(message):
		log.Printf("%s: command 47\n", branch.Name)
		fmt.Printf("%s: command 47\n", branch.Name)
	case REGEX_MSG_48.MatchString(message):
		log.Printf("%s: command 48\n", branch.Name)
		fmt.Printf("%s: command 48\n", branch.Name)
	case REGEX_MSG_50.MatchString(message):
		log.Printf("%s: command 50\n", branch.Name)
		fmt.Printf("%s: command 50\n", branch.Name)
	default:
		log.Printf("%s: unrecognized command for message %s\n", branch.Name, message)
		fmt.Printf("%s: unrecognized command for message %s\n", branch.Name, message)
	}
}

func getLatestGitMessage(workDir string, opts *RunnerOptions, branch Branch) (string, bool) {
	cfg := CommandConfig{
		Name: "Git log",
		Cmd:  "git",
		Args: []string{
			"log",
			"--no-decorate",
			"-n",
			"1",
		},
		WorkDir:     workDir,
		MaxAttempts: 1,
	}
	output, err := executeCommand(opts, cfg)
	if err != nil {
		log.Printf("Git diff error for branch %s: %v", branch.Name, err)
		return "", false
	}

	// Parse the output
	logMsg := output.Stdout
	cut := strings.Split(logMsg, "\n")
	if len(cut) < 6 {
		log.Printf("Unrecognized short(%d) log for %s:\n", len(cut), branch.Name)
		fmt.Printf("Unrecognized short(%d) log for %s:\n", len(cut), branch.Name)
		return "", false
	}

	// Find the author line, its not always at the same index.
	idx := 0
	var author string
	for idx < 5 {
		author = cut[idx]
		if regexp.MustCompile(`^Author: `).MatchString(author) {
			break
		}
		idx++
	}
	if !strings.Contains(author, "kcc-conductor-bot@google.com") {
		log.Printf("%s: unrecognized command by author %s\n", branch.Name, author)
		fmt.Printf("%s: unrecognized command by author %s\n", branch.Name, author)
		return "", false
	}

	// We believe the message will be at author index + 3
	if idx+3 >= len(cut) {
		log.Printf("%s: unrecognized command no message\n", branch.Name)
		fmt.Printf("%s: unrecognized command no message\n", branch.Name)
		return "", false
	}
	message := cut[idx+3]
	return message, true
}

func diffLastNCommits(opts *RunnerOptions, branch Branch) {
	close := setLoggingWriter(opts, branch)
	defer close()
	workDir := opts.branchRepoDir

	ctx := context.TODO()
	checkoutBranch(ctx, branch, workDir)

	// Run git diff command
	cfg := CommandConfig{
		Name: "Git diff",
		Cmd:  "git",
		Args: []string{
			"diff",
			"-r",
			fmt.Sprintf("HEAD~%d", opts.numCommits),
		},
		WorkDir:     workDir,
		MaxAttempts: 1,
	}
	output, err := executeCommand(opts, cfg)
	if err != nil {
		log.Printf("Git diff error for branch %s: %v", branch.Name, err)
		return
	}

	// Print the diff output
	log.Printf("Diff for branch %s (last %d commits):\n%s", branch.Name, opts.numCommits, output.Stdout)
	fmt.Printf("Diff for branch %s (last %d commits):\n%s", branch.Name, opts.numCommits, output.Stdout)
}

func revertLastNCommits(opts *RunnerOptions, branch Branch) {
	close := setLoggingWriter(opts, branch)
	defer close()
	workDir := opts.branchRepoDir
	ctx := context.TODO()
	checkoutBranch(ctx, branch, workDir)

	// First check for merge commits
	cfg := CommandConfig{
		Name: "Check for merge commits",
		Cmd:  "git",
		Args: []string{
			"log",
			fmt.Sprintf("-n%d", opts.numCommits),
			"--oneline",
		},
		WorkDir:     workDir,
		MaxAttempts: 1,
	}
	output, err := executeCommand(opts, cfg)
	if err != nil {
		log.Printf("Git log error for branch %s: %v", branch.Name, err)
		return
	}

	// Print the last N commits
	mergeCommit := false
	lines := strings.Split(strings.TrimSpace(output.Stdout), "\n")
	if len(lines) > 0 {
		for _, line := range lines {
			if strings.Contains(line, "Merge pull request") {
				mergeCommit = true
			}
		}
	} else {
		log.Printf("No commits found for branch %s", branch.Name)
		return
	}
	if mergeCommit {
		log.Printf("ERROR: Found merge commits in the last %d commits for branch %s:", opts.numCommits, branch.Name)
		fmt.Printf("\nERROR: Found merge commits in the last %d commits that would be reverted:\n%s\n", opts.numCommits, output.Stdout)
		return
	}

	// Show the diff
	diffLastNCommits(opts, branch)
	log.Printf("Last %d commits in branch %s:", opts.numCommits, branch.Name)
	for _, line := range lines {
		log.Printf("%s\n", line)
	}

	// Ask for final confirmation
	fmt.Printf("Are you sure you want to revert the last %d commits in branch %s? [y/N]: ", opts.numCommits, branch.Name)
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		log.Printf("Error reading input: %v", err)
		return
	}
	if response != "y" && response != "Y" {
		log.Printf("Skipping revert for branch %s", branch.Name)
		return
	}

	// Run git reset command
	cfg = CommandConfig{
		Name: "Git reset",
		Cmd:  "git",
		Args: []string{
			"reset",
			"--hard",
			fmt.Sprintf("HEAD~%d", opts.numCommits),
		},
		WorkDir:     workDir,
		MaxAttempts: 1,
	}
	output, err = executeCommand(opts, cfg)
	if err != nil {
		log.Printf("Git reset error for branch %s: %v", branch.Name, err)
		return
	}

	// Print the reset output
	log.Printf("Reset output for branch %s (last %d commits):\n%s", branch.Name, opts.numCommits, output.Stdout)
}
