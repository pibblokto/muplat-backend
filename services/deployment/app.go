package deployment

import (
	"errors"
	"fmt"
	"strings"

	"github.com/muplat/muplat-backend/pkg/k8s"
	"github.com/muplat/muplat-backend/repositories"
)

func CreateAppDeployment(
	deploymentName,
	projectName,
	deploymentType,
	username string,
	ac *AppConfig,
	db *repositories.Database,
	cc *k8s.ClusterConnection,
) error {
	if ac == nil {
		return errors.New("app config was not provided")
	}
	if *ac.External && ac.DomainName == "" {
		return errors.New("external flag was set but no domain was specified")
	}
	p, err := db.GetPorjectByName(projectName)
	if err != nil {
		return err
	}
	projectNamespace := p.Namespace

	u, err := db.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if p.Owner != username && !u.Admin {
		return errors.New("user lacks permissions to create deployment")
	}

	nameSuffix := k8s.GetNameSuffix(fmt.Sprintf("%s%s", deploymentName, projectName))
	resourceName := strings.ToLower(fmt.Sprintf("%s-%s", deploymentName, nameSuffix))
	resourceLabels := map[string]string{
		"name":         deploymentName,
		"created":      username,
		"project-name": projectName,
	}

	var secretName string
	if ac.EnvVars != nil {
		secretObject := cc.CreateSecretObject(
			resourceName,
			projectNamespace,
			resourceLabels,
			map[string]string{},
			ac.EnvVars,
		)
		err := cc.ApplySecret(secretObject)
		if err != nil {
			return err
		}
		secretName = resourceName
	} else {
		secretName = ""
	}

	// Deployment
	deploymentObject := cc.CreateDeploymentObject(
		resourceName,
		projectNamespace,
		resourceLabels,
		map[string]string{},
		ac.Repository,
		ac.Tag,
		ac.Port,
		secretName,
		1,
	)
	err = cc.ApplyDeployment(deploymentObject)
	if err != nil {

		deleteError := cc.DeleteSecret(resourceName, projectNamespace)
		if deleteError != nil {
			return fmt.Errorf("failed to create app %s and delete orphan resources", deploymentName)
		}

		return err
	}

	// Service
	serviceObject := cc.CreateServiceObject(
		resourceName,
		projectNamespace,
		resourceLabels,
		map[string]string{},
		resourceLabels["name"],
		ac.Port,
	)
	err = cc.ApplyService(serviceObject)
	if err != nil {
		deleteError := cc.DeleteDeployment(resourceName, projectNamespace)
		if deleteError != nil {
			return fmt.Errorf("failed to create app %s and delete orphan resources", deploymentName)
		}

		deleteError = cc.DeleteSecret(resourceName, projectNamespace)
		if deleteError != nil {
			return fmt.Errorf("failed to create app %s and delete orphan resources", deploymentName)
		}

		return err
	}

	// Ingress
	if *ac.External {
		ingressObject := cc.CreateIngressObject(
			resourceName,
			projectNamespace,
			resourceLabels,
			map[string]string{},
			ac.DomainName,
			resourceName,
			ac.Port,
		)
		err = cc.ApplyIngress(ingressObject)
		if err != nil {

			deleteError := cc.DeleteService(resourceName, projectNamespace)
			if deleteError != nil {
				return fmt.Errorf("failed to create app %s and delete orphan resources", deploymentName)
			}

			deleteError = cc.DeleteDeployment(resourceName, projectNamespace)
			if deleteError != nil {
				return fmt.Errorf("failed to create app %s and delete orphan resources", deploymentName)
			}

			deleteError = cc.DeleteSecret(resourceName, projectNamespace)
			if deleteError != nil {
				return fmt.Errorf("failed to create app %s and delete orphan resources", deploymentName)
			}
			return err
		}
	}

	err = db.SaveDeployment(deploymentName, projectName, deploymentType)
	if err != nil {
		return err
	}

	err = db.SaveAppConfig(
		deploymentName,
		projectName,
		ac.Repository,
		ac.Tag,
		getExternalUrl(ac),
		fmt.Sprintf("http://%s:%d", resourceName, ac.Port),
		string(ac.Tier),
		resourceName,
		ac.Port,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAppDeployment(
	deploymentName string,
	projectName string,
	username string,
	db *repositories.Database,
	cc *k8s.ClusterConnection,
) error {
	nameSuffix := k8s.GetNameSuffix(fmt.Sprintf("%s%s", deploymentName, projectName))
	resourceName := strings.ToLower(fmt.Sprintf("%s-%s", deploymentName, nameSuffix))

	p, err := db.GetPorjectByName(projectName)
	if err != nil {
		return err
	}
	projectNamespace := p.Namespace

	u, err := db.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if p.Owner != username && !u.Admin {
		return errors.New("user lacks permissions to delete deployment")
	}

	err = cc.DeleteIngress(resourceName, projectNamespace)
	if err != nil {
		return err
	}

	err = cc.DeleteService(resourceName, projectNamespace)
	if err != nil {
		return err
	}

	err = cc.DeleteDeployment(resourceName, projectNamespace)
	if err != nil {
		return err
	}

	err = cc.DeleteSecret(resourceName, projectNamespace)
	if err != nil {
		return err
	}

	err = db.DeleteAppConfig(deploymentName, projectName)
	if err != nil {
		return err
	}

	err = db.DeleteDeployment(deploymentName, projectName)
	if err != nil {
		return err
	}

	return nil
}

func getExternalUrl(ac *AppConfig) string {
	if *ac.External {
		return fmt.Sprintf("http://%s", ac.DomainName)
	}
	return ""
}
