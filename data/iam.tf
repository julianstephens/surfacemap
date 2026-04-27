resource "aws_iam_role" "lambda_execution_role" {
  name = "${var.project}-lambda-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "lambda_ddb_access" {
  statement {
    sid    = "DynamoDBAccess"
    effect = "Allow"
    actions = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem",
      "dynamodb:Scan",
    "dynamodb:Query", ]
    resources = [aws_dynamodb_table.this.arn, "${aws_dynamodb_table.this.arn}/*"]
  }
}

data "aws_iam_policy_document" "lambda_cw_access" {
  statement {
    sid    = "CloudWatchLogsAccess"
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
    "logs:PutLogEvents", ]
    resources = [aws_cloudwatch_log_group.lambda_log_group.arn, "${aws_cloudwatch_log_group.lambda_log_group.arn}:*"]
  }
}

resource "aws_iam_policy_attachment" "lambda_ddb_access_attachment" {
  name       = "${var.project}-lambda-ddb-access"
  policy_arn = aws_iam_policy.lambda_ddb_access.arn
  roles      = [aws_iam_role.lambda_execution_role.name]
}

resource "aws_iam_policy_attachment" "lambda_cw_access_attachment" {
  name       = "${var.project}-lambda-cw-access"
  policy_arn = aws_iam_policy.lambda_cw_access.arn
  roles      = [aws_iam_role.lambda_execution_role.name]
}