package cmd

import (
	"os"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/spf13/cobra"
)

type kubernetesDeployOptions struct {
	AdditionalParameters       []string `json:"additionalParameters,omitempty"`
	APIServer                  string   `json:"apiServer,omitempty"`
	AppTemplate                string   `json:"appTemplate,omitempty"`
	ChartPath                  string   `json:"chartPath,omitempty"`
	ContainerRegistryPassword  string   `json:"containerRegistryPassword,omitempty"`
	ContainerRegistryURL       string   `json:"containerRegistryUrl,omitempty"`
	ContainerRegistryUser      string   `json:"containerRegistryUser,omitempty"`
	CreateDockerRegistrySecret bool     `json:"createDockerRegistrySecret,omitempty"`
	DeploymentName             string   `json:"deploymentName,omitempty"`
	DeployTool                 string   `json:"deployTool,omitempty"`
	EnvVars                    []string `json:"envVars,omitempty"`
	HelmDeployWaitSeconds      int      `json:"helmDeployWaitSeconds,omitempty"`
	Image                      string   `json:"image,omitempty"`
	IngressHosts               []string `json:"ingressHosts,omitempty"`
	KubeConfig                 string   `json:"kubeConfig,omitempty"`
	KubeContext                string   `json:"kubeContext,omitempty"`
	Namespace                  string   `json:"namespace,omitempty"`
	TillerNamespace            string   `json:"tillerNamespace,omitempty"`
}

var myKubernetesDeployOptions kubernetesDeployOptions
var kubernetesDeployStepConfigJSON string

// KubernetesDeployCommand Deployment to Kubernetes test or production namespace within the specified Kubernetes cluster.
func KubernetesDeployCommand() *cobra.Command {
	metadata := kubernetesDeployMetadata()
	var createKubernetesDeployCmd = &cobra.Command{
		Use:   "kubernetesDeploy",
		Short: "Deployment to Kubernetes test or production namespace within the specified Kubernetes cluster.",
		Long: `Deployment to Kubernetes test or production namespace within the specified Kubernetes cluster.

!!! note "Deployment supports multiple deployment tools"
    Currently the following are supported:

    * [Helm](https://helm.sh/) command line tool and [Helm Charts](https://docs.helm.sh/developing_charts/#charts).
    * [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and ` + "`" + `kubectl apply` + "`" + ` command.

## Helm
Following helm command will be executed by default:

` + "`" + `` + "`" + `` + "`" + `
helm upgrade <deploymentName> <chartPath> --install --force --namespace <namespace> --wait --timeout <helmDeployWaitSeconds> --set "image.repository=<yourRegistry>/<yourImageName>,image.tag=<yourImageTag>,secret.dockerconfigjson=<dockerSecret>,ingress.hosts[0]=<ingressHosts[0]>,,ingress.hosts[1]=<ingressHosts[1]>,...
` + "`" + `` + "`" + `` + "`" + `

* ` + "`" + `yourRegistry` + "`" + ` will be retrieved from ` + "`" + `containerRegistryUrl` + "`" + `
* ` + "`" + `yourImageName` + "`" + `, ` + "`" + `yourImageTag` + "`" + ` will be retrieved from ` + "`" + `image` + "`" + `
* ` + "`" + `dockerSecret` + "`" + ` will be calculated with a call to ` + "`" + `kubectl create secret docker-registry regsecret --docker-server=<yourRegistry> --docker-username=<containerRegistryUser> --docker-password=<containerRegistryPassword> --dry-run=true --output=json'` + "`" + ``,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetStepName("kubernetesDeploy")
			log.SetVerbose(GeneralConfig.Verbose)
			return PrepareConfig(cmd, &metadata, "kubernetesDeploy", &myKubernetesDeployOptions, config.OpenPiperFile)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return kubernetesDeploy(myKubernetesDeployOptions)
		},
	}

	addKubernetesDeployFlags(createKubernetesDeployCmd)
	return createKubernetesDeployCmd
}

func addKubernetesDeployFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceVar(&myKubernetesDeployOptions.AdditionalParameters, "additionalParameters", []string{}, "Defines additional parameters for `helm install` or `kubectl apply` command.")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.APIServer, "apiServer", os.Getenv("PIPER_apiServer"), "Defines the Url of the API Server of the Kubernetes cluster.")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.AppTemplate, "appTemplate", os.Getenv("PIPER_appTemplate"), "Defines the filename for the kubernetes app template (e.g. k8s_apptemplate.yaml)")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.ChartPath, "chartPath", os.Getenv("PIPER_chartPath"), "Defines the chart path for deployments using helm.")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.ContainerRegistryPassword, "containerRegistryPassword", os.Getenv("PIPER_containerRegistryPassword"), "")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.ContainerRegistryURL, "containerRegistryUrl", os.Getenv("PIPER_containerRegistryUrl"), "http(s) url of the Container registry.")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.ContainerRegistryUser, "containerRegistryUser", os.Getenv("PIPER_containerRegistryUser"), "")
	cmd.Flags().BoolVar(&myKubernetesDeployOptions.CreateDockerRegistrySecret, "createDockerRegistrySecret", true, "Toggle to turn on Regsecret creation with a `deployTool:kubectl` deployment.")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.DeploymentName, "deploymentName", os.Getenv("PIPER_deploymentName"), "Defines the name of the deployment.")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.DeployTool, "deployTool", "kubectl", "Defines the tool which should be used for deployment.")
	cmd.Flags().StringSliceVar(&myKubernetesDeployOptions.EnvVars, "envVars", []string{"map[HELM_HOME:/home/piper/.helm KUBECONFIG:/home/piper/.kube/config]"}, "Environment variables which should be passed to HELM deployment.")
	cmd.Flags().IntVar(&myKubernetesDeployOptions.HelmDeployWaitSeconds, "helmDeployWaitSeconds", 300, "Number of seconds before helm deploy returns.")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.Image, "image", os.Getenv("PIPER_image"), "Full name of the image to be deployed.")
	cmd.Flags().StringSliceVar(&myKubernetesDeployOptions.IngressHosts, "ingressHosts", []string{}, "List of ingress hosts to be exposed via helm deployment.")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.KubeConfig, "kubeConfig", os.Getenv("PIPER_kubeConfig"), "Defines the path to the `kubeconfig` file.")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.KubeContext, "kubeContext", os.Getenv("PIPER_kubeContext"), "Defines the context to use from the `kubeconfig` file.")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.Namespace, "namespace", os.Getenv("PIPER_namespace"), "Defines the target Kubernetes namespace for the deployment.")
	cmd.Flags().StringVar(&myKubernetesDeployOptions.TillerNamespace, "tillerNamespace", os.Getenv("PIPER_tillerNamespace"), "Defines optional tiller namespace for deployments using helm.")

	cmd.MarkFlagRequired("apiServer")
	cmd.MarkFlagRequired("chartPath")
	cmd.MarkFlagRequired("deploymentName")
	cmd.MarkFlagRequired("deployTool")
	cmd.MarkFlagRequired("image")
}

// retrieve step metadata
func kubernetesDeployMetadata() config.StepData {
	var theMetaData = config.StepData{
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Parameters: []config.StepParameters{
					{
						Name:      "additionalParameters",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "[]string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "helmDeploymentParameters"}},
					},
					{
						Name:      "apiServer",
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{{Name: "k8sAPIServer"}},
					},
					{
						Name:      "appTemplate",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "k8sAppTemplate"}},
					},
					{
						Name:      "chartPath",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{{Name: "helmChartPath"}},
					},
					{
						Name:      "containerRegistryPassword",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name:      "containerRegistryUrl",
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "dockerRegistryUrl"}},
					},
					{
						Name:      "containerRegistryUser",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name:      "createDockerRegistrySecret",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "bool",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name:      "deploymentName",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{{Name: "helmDeploymentName"}},
					},
					{
						Name:      "deployTool",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{},
					},
					{
						Name:      "envVars",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "[]string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "helmEnvVars"}},
					},
					{
						Name:      "helmDeployWaitSeconds",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "int",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name:      "image",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{{Name: "deployImage"}},
					},
					{
						Name:      "ingressHosts",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "[]string",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name:      "kubeConfig",
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name:      "kubeContext",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name:      "namespace",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "helmDeploymentNamespace"}, {Name: "k8sDeploymentNamespace"}},
					},
					{
						Name:      "tillerNamespace",
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "helmTillerNamespace"}},
					},
				},
			},
		},
	}
	return theMetaData
}