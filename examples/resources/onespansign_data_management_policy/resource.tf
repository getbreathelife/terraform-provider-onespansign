resource "onespansign_data_management_policy" "example" {
  transaction_retention {
    draft                     = 0
    sent                      = 0
    completed                 = 0
    archived                  = 0
    declined                  = 0
    opted_out                 = 0
    expired                   = 0
    lifetime_total            = 120
    lifetime_until_completion = 120
    include_sent              = false
  }
}