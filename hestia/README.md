# README

```bash
bundle exec rails generate model Company name alias_name intro:text legal_person tel website email address established_at:datetime unified_social_credit_code logo deleted_at:datetime

bundle exec rails generate model User username password mobile email avatar note ua ip os browser location:point company_id:integer first_login_at:datetime last_active_at:datetime deleted_at:datetime

bundle exec rails generate model ChatMessage from:integer to:integer content:text read_at:datetime deleted_at:datetime

bundle exec rails g model App uuid name company_id:integer intro:text domains deleted_at:datetime
```
