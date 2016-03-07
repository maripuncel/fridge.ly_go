namespace :db do
  namespace :prod do
    task :migrate do
      sh "migrate -url $(heroku run 'echo $DATABASE_URL') -path ./migrations up"
    end
  end

  namespace :local do
    database_name = "fridgly_devel"
    task :setup do
        sh "psql -h localhost postgres -c 'CREATE DATABASE fridgly_devel'"
    end

    task :migrate do
      sh "migrate -url postgres://localhost/#{database_name}?sslmode=disable -path ./migrations up"
    end
  end
end
