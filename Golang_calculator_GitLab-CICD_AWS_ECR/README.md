### Create and test Golang calculator with the GitLab-CI/CD, Podman executor and deploy an image to the AWS ECR

The new project in the GitLab.com is go-calculator



#### Create and push Podman's image to the GitLab.com

Create image from Containerfiles 


`localhost$ podman build -t calc-test -f Containerfile.test`

```
localhost$ podman images | grep 'calc-*'
localhost/calc-test                                  latest               75d343c5722d  About 20 minutes ago 11.9 MB

```

Create project access token in GitLab.com __Settings__ -> __Access tokens__ -> __Add new token__ with __write_registry__ and 
__read_registry__ permissions. Save username __admin__ and password.

Login podman with the GitLab.com container registry

```
localhost$ podman login -u admin -p ****** --authfile ~/.config/containers/go-calculator/auth.json registry.gitlab.com/user_name/go-calculator
   Login Succeeded!

```

Push test image to the GitLab.com container registry


```

localhost$ podman push --authfile ~/.config/containers/go-calculator/auth.json localhost/calc-test:latest
 registry.gitlab.com/user_name/go-calculator/calc_test
Getting image source signatures
Copying blob ae977f3f1973 done  
Copying blob 78561cef0761 done  
Copying config 75d343c572 done  
Writing manifest to image destination

```

#### AWS needed steps

1. Create IAM __admin__ user, set AdministratorAccess and AmazonECS_FullAccess permissions and save credentions of this user

2. Create private repository in the ECR with __gitlab-repo__ name.


#### Create some variables in the GitLab.com

Copy content from saved IAM __admin__ user credentials aws_access_key_id to the $AWS_ACCESS_KEY_ID and 
aws_secret_access_key to the $AWS_SECRET_ACCESS_KEY.

Both variables are Masked and Expanded.

Assign $AWS_DEFAULT_REGION as IAM user working regioni us-east-1 for example.

Assign $AWS_ECR_URI to <Account ID>.dkr.ecr.<region>.amazonaws.com

These variables can be Visible


#### GitLab.com pipeline .pre job

This job to establish aws cli connection from GitLab.com to the Amazon ECR

```

get_aws_ecr_creds:
  stage: .pre
  image:
# official image
    name: amazon/aws-cli:2.27.46
    entrypoint: [""]
  script:
# try to connect to the AWS ECR, if credentions are valid save its 
   - aws ecr get-login-password --region $AWS_DEFAULT_REGION > aws_ecr_creds.txt
  artifacts:
    paths:
      - aws_ecr_creds.txt
```


when the job executed

```

...
$ aws ecr get-login-password --region $AWS_DEFAULT_REGION > aws_ecr_creds.txt
Uploading artifacts for successful job
00:02
Uploading artifacts...
aws_ecr_creds.txt: found 1 matching artifact files and directories 
Uploading artifacts as "archive" to coordinator... 201 Created  correlation_id=0b1bd57c80632b4ae988823ba70e982e id=10540766029 responseStatus=201 Created token=68_NasrTs
Cleaning up project directory and file based variables
00:01
Job succeeded
```

#### GitLab.com pipeline test_go_calc job

Create __test_go_calc__ job inside __test__ stage to test go-calculator application with Podman executor 

```

test_go_calc:
  stage: test
# official image from quay.io
  image: quay.io/containers/podman
  script:
# to get access to the GitLab.com registry
    - podman login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - podman run -dt --name go-calc-test $CI_REGISTRY_IMAGE/calc_test
# container run and exited
    - podman ps -a
    - podman logs go-calc-test
``` 

inside the test_go_calc job description

```

...
$ podman logs go-calc-test
=== RUN   TestAdd
=== PAUSE TestAdd
=== RUN   TestSubtract
=== PAUSE TestSubtract
=== RUN   TestMultiply
=== PAUSE TestMultiply
=== RUN   TestDivide
=== PAUSE TestDivide
=== RUN   TestSqrt
=== PAUSE TestSqrt
=== CONT  TestAdd
--- PASS: TestAdd (0.00s)
=== CONT  TestSqrt
--- PASS: TestSqrt (0.00s)
=== CONT  TestDivide
--- PASS: TestDivide (0.00s)
=== CONT  TestMultiply
--- PASS: TestMultiply (0.00s)
=== CONT  TestSubtract
--- PASS: TestSubtract (0.00s)
PASS
Cleaning up project directory and file based variables
00:00
Job succeeded

```

#### GitLab.com pipeline deploy_img_to_aws job

This job builds image from Containerfile.build and deploy it to the AWS ECR

```

deploy_img_to_aws:
  stage: deploy
  image: quay.io/containers/podman
  script:
# build image with podman
    - podman build -t go-calc-build -f Containerfile.build .
# podman login to the AWS ECR with artifacts credentials from .pre job
    - cat aws_ecr_creds.txt | podman login --username AWS  --password-stdin  $AWS_ECR_URI
# push image to the private AWS repository
    - podman push go-calc-build $AWS_ECR_URI/gitlab-repo
```

when the job completed

```

$ podman push go-calc-build <Account ID>.dkr.ecr.us-east-1.amazonaws.com/gitlab-repo
Getting image source signatures
Copying blob sha256:30771c8b778f38158ae488c8b02162482804b9c0292e85ff029ae336af264c86
Copying blob sha256:fd2758d7a50e2b78d275ee7d1c218489f2439084449d895fa17eede6c61ab2c4
Copying config sha256:6afb59e3d9642d319c31d5e1e597dc4b180aea2905f80ea6bb1b6bca49a64488
Writing manifest to image destination
Cleaning up project directory and file based variables
00:01
Job succeeded
```

#### GitLab.com pipeline run_img_from_aws job - create and run container from AWS ECR image

This job downloads image from AWS ECR runs container and gets it logs

```
run_img_from_aws:
  stage: run
  image: quay.io/containers/podman
  script:
    - cat aws_ecr_creds.txt | podman login --username AWS --password-stdin $AWS_ECR_URI
    - podman run -dt --name go-calc $AWS_ECR_URI/gitlab-repo:latest
# Container is exited so --all arguments is needed
    - podman ps -a
    - podman logs go-calc
```

when the job finished

```

$ cat aws_ecr_creds.txt | podman login --username AWS --password-stdin $AWS_ECR_URI
Login Succeeded!
$ podman run -dt --name go-calc $AWS_ECR_URI/gitlab-repo:latest
Trying to pull 371439860319.dkr.ecr.us-east-1.amazonaws.com/gitlab-repo:latest...
Getting image source signatures
Copying blob sha256:4b00fdf0472518f0ae444cae5ba0c697bd75af219272889afa1b893dab58cb0d
Copying blob sha256:0b7bc71b010c407d231087ab07b6fa4b2ea8a8262956a9ff299996d6fabea4f0
Copying config sha256:41d0ef58a165fdfa977a7a67cc5fa2dad69d067d3a41731b8386adb5bb29dcc2
Writing manifest to image destination
4a712f44e47bfbd2dcd8923186e8d3156dacec8e8f41ff6e0ae40f72a00fdbc9
$ podman ps -a
CONTAINER ID  IMAGE                                                            COMMAND     CREATED                 STATUS                 PORTS       NAMES
4a712f44e47b  371439860319.dkr.ecr.us-east-1.amazonaws.com/gitlab-repo:latest  ./calc      Less than a second ago  Up Less than a second              go-calc
$ podman logs go-calc
Add 4 to 5 is 9
Substract 5 from 4 is -1
Multiply 4 by 5 is 20
Divide 4 by 5 is 0.8
Sqrt of 5 is 2.23606797749979
$ podman images
REPOSITORY                                                TAG         IMAGE ID      CREATED         SIZE
371439860319.dkr.ecr.us-east-1.amazonaws.com/gitlab-repo  latest      41d0ef58a165  30 seconds ago  10.9 MB
Cleaning up project directory and file based variables
00:00
Job succeeded
```

