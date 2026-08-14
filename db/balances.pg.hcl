table "balances" {
  schema = schema.public

  column "id" {
    type = varchar(36)
    null = false
  }

  column "account_id" {
    type = varchar(36)
    null = false
  }

  column "amount" {
    type = bigint
    null = false
  }

  column "currency" {
    type = varchar(5)
    null = false
  }

  column "created_at" {
    type = timestamptz
    null = false
  }

  column "updated_at" {
    type = timestamptz
    null = false
  }
  
  primary_key {
    columns = [ column.id ]
  }

  foreign_key "fk_account_id" {
    columns = [column.account_id]
    ref_columns = [table.accounts.column.id]
  }

  index "idx_balance_account_number_currency" {
    columns = [column.account_id, column.currency]
    unique  = true
  }

}