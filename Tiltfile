k8s_yaml('deployment.yaml')
k8s_yaml('namespace.yaml')
k8s_yaml('rbac.yaml')

# Provide some test data
k8s_yaml('testdata/projects/testdata.yaml')

custom_build(
  'telliottio/front',
  'make build',
  [
    'pkg',
    'internal',
    'cmd',
    'public',
    'views',
    'deployment.yaml',
    'Dockerfile',
    'ingress.yaml',
    'Makefile',
    'go.mod'
  ],
  tag="latest"
)

k8s_resource('front', port_forwards='8080:80')