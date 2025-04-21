package deployments

//if input.AppConfing == nil {
//	c.JSON(http.StatusBadRequest, "App config was not provided")
//	return
//}
//deploymentObject := k8s.CreateDeploymentObject(
//	deploymentName,
//	deploymentNamespace,
//	deploymentLabels,
//	map[string]string{},
//	string(input.AppConfing.Tier),
//	input.AppConfing.Repository,
//	input.AppConfing.Tag,
//	input.AppConfing.Port,
//	"",
//)
//err = k8s.ApplyDeployment(clientset, deploymentObject)
//if err != nil {
//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//	return
//}
