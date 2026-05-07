package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

// ECRClient defines the subset of the ECR API used by FetchECRRepository.
type ECRClient interface {
	DescribeRepositories(ctx context.Context, params *ecr.DescribeRepositoriesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeRepositoriesOutput, error)
}

// FetchECRRepository retrieves attributes of an ECR repository by name.
func FetchECRRepository(ctx context.Context, client ECRClient, repositoryName string) (map[string]string, error) {
	if repositoryName == "" {
		return nil, fmt.Errorf("ecr: repository name must not be empty")
	}

	out, err := client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{repositoryName},
	})
	if err != nil {
		return nil, fmt.Errorf("ecr: describe repositories: %w", err)
	}

	if len(out.Repositories) == 0 {
		return nil, fmt.Errorf("ecr: repository %q not found", repositoryName)
	}

	repo := out.Repositories[0]
	attrs := map[string]string{
		"repository_name": aws.ToString(repo.RepositoryName),
		"repository_uri":  aws.ToString(repo.RepositoryUri),
		"registry_id":     aws.ToString(repo.RegistryId),
	}

	if repo.ImageScanningConfiguration != nil {
		if repo.ImageScanningConfiguration.ScanOnPush {
			attrs["image_scanning_scan_on_push"] = "true"
		} else {
			attrs["image_scanning_scan_on_push"] = "false"
		}
	}

	if string(repo.ImageTagMutability) != "" {
		attrs["image_tag_mutability"] = string(repo.ImageTagMutability)
	}

	return attrs, nil
}
