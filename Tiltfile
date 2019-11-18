# Deploy: tell Tilt what YAML to deploy
k8s_yaml('deployment.yaml')

custom_build(
  'telliottio/front',
  'make build',
  ['cmd', 'public', 'views', 'deployment.yaml', 'Dockerfile', 'ingress.yaml', 'Makefile', 'go.mod"'],
  tag="latest"
)

k8s_resource('front', port_forwards='8080:80')