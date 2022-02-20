# Waypoint GitLab Plugin

Plugin for Hashicorp Waypoint to push artifacts to a GitLab registry

## Usage

### Inputs

- `name` = name of the package being pushed
- `version` = version # of the package
- `project_id` = GitLab project ID
- `file_name` = name of the file to upload to GitLab
- `token` = GitLab API token for pushing to registry

### Example

```hcl
project = "hello-go"

plugin "gitlab" {
  type {
    registry = true
    mapper = true
  }
}

app "hello-go" {
  build {
    use "go" {
      source      = "./"
      output_name = "hello"
    }
    registry {
      use "gitlab" {
        name       = "hello-app"
        version    = "1.0.0"
        project_id = 123456789
        file_name  = "hello"
        token      = "MY_GITLAB_TOKEN"
      }
    }
  }

  deploy {
    use "exec" {
      command = ["echo", "Build", "complete!"]
    }
  }
}
```
