# frozen_string_literal: true

task default: %w[clean]

task :clean do
  p 'clean task'
end

task :fmt do
  sh 'go fmt ./...'
end

DB_URL = 'postgres://yang:@localhost:5432/hestia_development?sslmode=disable'

namespace :g do
  task :model do
    sh "xorm reverse postgres '#{DB_URL}' $PWD/models/templates/goxorm $PWD/models"
    sh 'rm -f $PWD/models/ar_internal_metadata.go'
    sh 'rm -f $PWD/models/schema_migrations.go'
    sh 'rm -f $PWD/models/spatial_ref_sys.go'
  end
end
