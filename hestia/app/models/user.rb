class User < ApplicationRecord
  has_secure_password
  validates :role, inclusion: { in: %w[admin member customer], message: '%{value} is not a valid role' }

  belongs_to :company
end
