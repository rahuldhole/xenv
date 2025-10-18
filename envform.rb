#!/usr/bin/env ruby
# EnvForm - Interactive environment variable configuration tool
# Usage: ruby envform.rb env.form

require 'io/console'

if ARGV.length < 1
  puts "Usage: #{$0} <form-file>"
  exit 1
end

FORM_FILE = ARGV[0]
OUTPUT_FILE = FORM_FILE.sub(/\.form$/, '')

unless File.exist?(FORM_FILE)
  puts "Error: Form file '#{FORM_FILE}' not found."
  exit 1
end

# Function to read existing value for a key
def get_existing_value(key, default_value, output_file)
  return default_value unless File.exist?(output_file)
  
  File.readlines(output_file).each do |line|
    if line =~ /^#{Regexp.escape(key)}=(.*)$/
      return $1.chomp
    end
  end
  
  default_value
end

# Function to handle text input
def prompt_input(label, default_value)
  print "#{label}"
  print " [#{default_value}]" unless default_value.empty?
  print ": "
  
  input = STDIN.gets.chomp
  input.empty? ? default_value : input
end

# Function to handle password input
def prompt_password(label, default_value)
  if default_value && !default_value.empty?
    print "#{label} [press enter to keep current value]: "
  else
    print "#{label}: "
  end
  
  input = STDIN.noecho(&:gets).chomp
  puts
  
  input.empty? ? default_value : input
end

# Function to handle select input
def prompt_select(label, options, default_value)
  puts "#{label}:"
  
  options.each_with_index do |opt, idx|
    if opt == default_value
      puts "  #{idx + 1}) #{opt} (current)"
    else
      puts "  #{idx + 1}) #{opt}"
    end
  end
  
  print "Select [1-#{options.length}]: "
  selection = STDIN.gets.chomp
  
  return default_value if selection.empty?
  
  idx = selection.to_i - 1
  if idx >= 0 && idx < options.length
    options[idx]
  else
    default_value
  end
end

# Function to handle checkbox input
def prompt_checkbox(label, default_value)
  # Normalize default value
  default = case default_value.to_s.downcase
  when 'true', 'yes', 'y', '1', 'on'
    true
  else
    false
  end
  
  if default
    print "#{label} [Y/n]: "
  else
    print "#{label} [y/N]: "
  end
  
  response = STDIN.gets.chomp.downcase
  
  case response
  when 'y', 'yes', 'true'
    'true'
  when 'n', 'no', 'false'
    'false'
  when ''
    default ? 'true' : 'false'
  else
    default ? 'true' : 'false'
  end
end

# Process form and create new output
puts "Interactive configuration for #{File.basename(FORM_FILE)}"
puts "-" * 50

output_lines = []
last_field_type = nil
last_field_label = nil
last_field_options = nil

File.readlines(FORM_FILE).each do |line|
  line = line.chomp
  
  # Check if this is a comment line with field directives
  if line =~ /^#/
    if line =~ /@input/
      last_field_type = :input
      last_field_label = line[/label="([^"]*)"/, 1]
    elsif line =~ /@password/
      last_field_type = :password
      last_field_label = line[/label="([^"]*)"/, 1]
    elsif line =~ /@select/
      last_field_type = :select
      last_field_label = line[/label="([^"]*)"/, 1]
      last_field_options = line[/options=([^\s]*)/, 1]&.split(',') || []
    elsif line =~ /@checkbox/
      last_field_type = :checkbox
      last_field_label = line[/label="([^"]*)"/, 1]
    else
      last_field_type = nil
      last_field_label = nil
      last_field_options = nil
    end
    
    output_lines << line
  elsif line =~ /^([^=]+)=(.*)$/
    # This is a key=value line
    key = $1
    default_value = $2
    
    # Get existing value
    existing_value = get_existing_value(key, default_value, OUTPUT_FILE)
    
    # If we have a field type, use it for prompting
    if last_field_type
      # Use key as label if no label specified
      label = last_field_label || key
      
      new_value = case last_field_type
      when :input
        prompt_input(label, existing_value)
      when :password
        prompt_password(label, existing_value)
      when :select
        prompt_select(label, last_field_options, existing_value)
      when :checkbox
        prompt_checkbox(label, existing_value)
      else
        existing_value
      end
      
      # Reset field info
      last_field_type = nil
      last_field_label = nil
      last_field_options = nil
    else
      # No field type, just use the existing/default value
      new_value = existing_value
    end
    
    # Write updated key=value
    output_lines << "#{key}=#{new_value}"
  else
    # Other lines (blank lines etc.)
    output_lines << line
  end
end

# Write to output file
File.write(OUTPUT_FILE, output_lines.join("\n") + "\n")

puts "-" * 50
puts "Configuration saved to #{OUTPUT_FILE}"
