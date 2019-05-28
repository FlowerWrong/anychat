# frozen_string_literal: true

class LoginLog < ApplicationRecord
  belongs_to :user
end
