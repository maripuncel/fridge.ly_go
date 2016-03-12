namespace :db do
  namespace :prod do
    task :migrate do
      sh "migrate -url $(heroku run 'echo $DATABASE_URL') -path ./migrations up"
    end
  end

  namespace :local do
    devel_database = "fridgely_devel"
    test_database = "fridgely_test"
    task :setup do
        sh "psql -h localhost postgres -c 'CREATE DATABASE #{devel_database}'"
        sh "psql -h localhost postgres -c 'CREATE DATABASE #{test_database}'"
    end

    task :migrate do
      sh "migrate -url postgres://localhost/#{devel_database}?sslmode=disable -path ./migrations up"
      sh "migrate -url postgres://localhost/#{test_database}?sslmode=disable -path ./migrations up"
    end
  end
end
