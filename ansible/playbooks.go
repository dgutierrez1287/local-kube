package ansible

type Playbook struct {
  Name string         `json:"name" yaml:"name"`
  Hosts string        `json:"hosts" yaml:"hosts"`
  Become bool         `json:"become" yaml:"become"`
  BecomeUser string   `json:"become_user" yaml:"become_user"`
  VarsFiles []string  `json:"vars_files" yaml:"vars_files"`
  Roles []string      `json:"roles" yaml:"roles"`
}

