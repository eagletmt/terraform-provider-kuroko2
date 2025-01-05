FROM ruby:3.3

WORKDIR /app
COPY Gemfile Gemfile.lock ./
ENV BUNDLE_PATH=/gems BUNDLE_DEPLOYMENT=1 BUNDLE_JOBS=8
RUN bundle install
COPY ./ ./
CMD ["./run.sh"]
