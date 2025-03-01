gcloud auth login
gcloud config get-value project
gcloud config set project local-dev-353516
gcloud auth application-default login
brew install pulumi/tap/pulumi

pip3 install pulumi
pip3 install pulumi_gcp

pulumi new gcp-go
pulumi config set gcp:project local-dev-353516
pulumi up -s dev
pulumi destroy -s dev
gcloud compute instances list

