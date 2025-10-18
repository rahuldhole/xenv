FROM ruby:3.2

# Set working directory
WORKDIR /app

# Copy Gemfile and install dependencies
COPY Gemfile .
RUN bundle install

# Copy the rest of the app
COPY . .

# Set the entrypoint to run xenv.rb with a form file argument
# (You can override the CMD when running the container)
ENTRYPOINT ["ruby", "xenv.rb"]
CMD ["env.xenv"]