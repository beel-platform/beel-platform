#!/bin/bash

case ${1} in
  *.cfg) CONFIG_FILE=${1} ;;
  *) CONFIG_FILE=./blue.cfg ;;
esac
case `uname -s` in
  'Linux') CONFIG_FILE=$(readlink -e ${CONFIG_FILE}) ;;
  'Darwin') CONFIG_FILE=$(stat -f ${CONFIG_FILE}) ;;
  *) echo "Operating System not supported"; exit 1 ;;
esac
if [ ! -f ${CONFIG_FILE} ]; then
  echo "Provide a valid configuration file."
  exit 1
fi
source ${CONFIG_FILE}
IMG_NAME="bluespark/${name}"
IMG_TAG=${version}
IMG_BASE=${image_base}
case `uname -s` in
  'Linux')
    if [[ ! -d `readlink -e ${code_base} 2>/dev/null` ]]; then
      echo "[ERROR] Code base location ($code_base) not found."
      exit 1
    fi
    code_base=$(readlink -e ${code_base})
    ;;
  'Darwin') code_base=$(stat -f ${code_base}) ;;
esac
if [ ! -d ${code_base} ]; then
  echo "Code base location ($code_base) not found"
  exit 1
fi
DKR_VOLS=("${code_base}:${volume}")
DKR_PORTS=${ports}
HAB_PKGS=${packages}
AWS_REG_ID=${ecr_id}
AWS_REGION=${ecr_region}

function upload_image ()
{
  echo "Do you want to upload the images? (y/n): "
  read input_upload_image
  case ${input_upload_image} in
    'n') exit 0 ;;
    'y')
      docker tag ${IMG_NAME}:latest ${AWS_REG_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${IMG_NAME}:latest && \
      docker tag ${IMG_NAME}:latest ${AWS_REG_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${IMG_NAME}:${IMG_TAG}
      ( docker push ${AWS_REG_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${IMG_NAME}:latest || \
        aws ecr create-repository --repository-name ${IMG_NAME} && \
        docker push ${AWS_REG_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${IMG_NAME}:latest ) && \
      docker push ${AWS_REG_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${IMG_NAME}:${IMG_TAG} && \
      echo "Image ${IMG_NAME}:${IMG_TAG} has been uploaded."
      ;;
    *) echo "Not a valid answer"; upload_image ;;
  esac
}

function build_image ()
{
  echo "Do you want to create the image? (y/n):"
  read input_build_image
  case ${input_build_image} in
    'n') exit 0 ;;
    'y')
      cat <<-EOF > ./Dockerfile
FROM ${AWS_REG_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${IMG_BASE}:latest
$(for i in ${HAB_PKGS}; do echo "RUN hab pkg install ${i} && hab sup load ${i}"; done)
$(for i in ${DKR_PORTS}; do echo "EXPOSE ${i}"; done)
ENTRYPOINT ["hab"]
CMD ["sup", "run"]
EOF
      docker build -t ${IMG_NAME}:latest . && rm -f ./Dockerfile
      upload_image
      ;;
    *) echo "Not a valid answer"; build_image ;;
  esac
}

function aws_login ()
{
  aws --v 2>/dev/null
  if [[ $? == 0 ]]; then
    (aws ecr get-login --no-include-email --region ${AWS_REGION} || echo "Cannot connect to AWS ECR."; exit 1) | bash
  else
    echo "[ERROR] AWS Command Line Interface is not installed"
    exit 1
  fi
}

# Validate if Docker is running
if [[ ! `docker -v 2>/dev/null` ]]; then
  echo "[ERROR] Docker is not installed"
  exit 1
fi

if [[ `docker images 2>&1 | grep 'Cannot connect' 2>/dev/null` ]]; then
  echo "[ERROR] Docker is not running"
  exit 1
elif [[ `docker images --format "{{.Repository}}:{{.Tag}}" | grep ${IMG_NAME}:${IMG_TAG} 2>/dev/null` ]]; then
  echo "Docker image ${IMG_NAME} found locally."
  aws_login
  if [[ ! `aws ecr list-images --registry-id ${AWS_REG_ID} --repository-name ${IMG_NAME} | grep ${IMG_TAG} 2>/dev/null` ]]; then
    echo "Docker image ${IMG_NAME} not found in the registry."
    upload_image
  fi
else
  echo "Docker image not found locally, checking in the registry"
  aws_login
  if [[ `aws ecr list-images --registry-id ${AWS_REG_ID} --repository-name ${IMG_NAME} | grep ${IMG_TAG} 2>/dev/null` ]]; then
    echo "Docker image found in registry."
    docker pull ${AWS_REG_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${IMG_NAME}:${IMG_TAG}
  else
    echo "Docker image not found"
    build_image
  fi
fi

if [[ ! `docker ps -a --format "{{.Image}}" | grep ${IMG_NAME}:${IMG_TAG}` ]]; then
  docker run --rm -td -p $(for i in ${DKR_PORTS}; do echo "${i}:${i}"; done) \
  -v $(for i in ${DKR_VOLS[@]}; do echo "${i}"; done) \
  ${AWS_REG_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${IMG_NAME}:${IMG_TAG}
  echo "Running: docker run -td -p $(for i in ${DKR_PORTS}; do echo "${i}:${i}"; done) \
-v $(for i in ${DKR_VOLS[@]}; do echo "${i}"; done) \
${AWS_REG_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${IMG_NAME}:${IMG_TAG}"
else
  echo "Container with image ${IMG_NAME}:${IMG_TAG} already running."
  exit 0
fi

# docker exec -ti <CONTAINER> /bin/bash
