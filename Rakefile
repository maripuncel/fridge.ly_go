namespace :migrate do
  task :prod do
    sh "migrate -url $(heroku run 'echo $DATABASE_URL') -path ./migrations up"
  end

  task :local do
    database_name = "fridgly_devel"
    sh "migrate -url postgres://localhost/#{database_name}?sslmode=disable -path ./migrations up"
  end
end

task :setup do
    sh "psql -h localhost postgres -c 'CREATE DATABASE fridgly_devel'"
end
