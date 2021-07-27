resource "aws_db_parameter_group" "toebeans" {
  name   = "toebeans"
  family = "mysql5.7"

#   parameter {
#     name  = "character_set_database"
#     value = "utf8mb4"
#   }
#   parameter {
#     name  = "collation_server"
#     value = "utf8mb4_bin"
#   }
#   parameter {
#     name  = "default_time_zone"
#     value = "Asia/Tokyo"
#   }
#   parameter {
#     name  = "log_timestamps"
#     value = "SYSTEM"
#   }
#   parameter {
#     name  = "default_authentication_plugin"
#     value = "mysql_native_password"
#   }
#   parameter {
#     name  = "log_error"
#     value = "/var/log/mysql/mysql-error.log"
#   }
#   parameter {
#     name  = "slow_query_log"
#     value = 1
#   }
#   parameter {
#     name  = "slow_query_log_file"
#     value = "/rdsdbdata/log/slowquery/mysql-slowquery.log"
#   }
#   parameter {
#     name  = "long_query_time"
#     value = 5.0
#   }
#   parameter {
#     name  = "log_queries_not_using_indexes"
#     value = 0
#   }
#   parameter {
#     name  = "general_log"
#     value = 1
#   }
#   parameter {
#     name  = "general_log_file"
#     value = "/var/log/mysql/mysql-query.log"
#   }
}

resource "aws_db_option_group" "toebeans" {
  name                 = "toebeans"
  engine_name          = "mysql"
  major_engine_version = "5.7"

  option {
    option_name = "MARIADB_AUDIT_PLUGIN"
  }
}

resource "aws_db_subnet_group" "toebeans" {
  name       = "toebeans"
  subnet_ids = [aws_subnet.private_0.id, aws_subnet.private_1.id]
}

# resource "aws_db_instance" "toebeans" {
#   identifier                 = "toebeans"
#   engine                     = "mysql"
#   engine_version             = "5.7.25"
#   instance_class             = "db.t3.small"
#   allocated_storage          = 20
#   max_allocated_storage      = 100
#   storage_type               = "gp2"
#   storage_encrypted          = true
#   kms_key_id                 = aws_kms_key.toebeans.arn
#   username                   = "admin"
#   password                   = var.rds_password
#   multi_az                   = true
#   publicly_accessible        = false
#   backup_window              = "09:10-09:40"
#   backup_retention_period    = 30
#   maintenance_window         = "mon:10:10-mon:10:40"
#   auto_minor_version_upgrade = false
#   deletion_protection        = false // CAUTION
#   skip_final_snapshot        = true // CAUTION
#   port                       = 3306
#   apply_immediately          = false
#   vpc_security_group_ids     = [module.mysql_sg.security_group_id]
#   parameter_group_name       = aws_db_parameter_group.toebeans.name
#   option_group_name          = aws_db_option_group.toebeans.name
#   db_subnet_group_name       = aws_db_subnet_group.toebeans.name

#   lifecycle {
#     ignore_changes = [password]
#   }
# }

module "mysql_sg" {
  source      = "../../modules/security_group"
  name        = "mysql-sg"
  vpc_id      = aws_vpc.toebeans.id
  port        = 3306
  cidr_blocks = [aws_vpc.toebeans.cidr_block]
}
