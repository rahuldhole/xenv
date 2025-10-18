require 'rake'
require 'rspec/core/rake_task'

# Define a RSpec task
RSpec::Core::RakeTask.new(:spec) do |t|
  t.pattern = 'spec/**/*_spec.rb'
end

# Define a task to run all tests
task default: :spec

# Define additional tasks as needed
desc 'Run the application'
task :run do
  ruby 'lib/my_ruby_project.rb'
end

desc 'Clean up temporary files'
task :clean do
  rm_rf 'tmp/*'
end