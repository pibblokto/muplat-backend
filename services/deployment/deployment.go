package deployment

import (
	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/repositories"
)

func GetDeployment(deploymentName, projectName, callerUsername string, db *repositories.Database) (*gin.H, error) {
	d, err := db.GetDeployment(deploymentName, projectName)
	if err != nil {
		return nil, err
	}

	_, err = db.GetUserByUsername(callerUsername)
	if err != nil {
		return nil, err
	}

}
