env "local" {
  src = "./ent/schema"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ .Name }}"
    }
  }
}
