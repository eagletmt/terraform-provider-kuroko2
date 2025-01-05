#!/bin/bash
set -ex

bin/rails db:migrate
bin/rails s -p 3000 -b 0.0.0.0
