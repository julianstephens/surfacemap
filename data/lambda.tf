resource "aws_lambda_function" "this" {
  function_name = var.lambda.function_name
  role          = var.lambda.role_arn
  handler       = var.lambda.handler
  runtime       = var.lambda.runtime
  filename      = var.lambda.filename

  logging_config {
    log_format = "json"
  }

  tags = local.tags
}

resource "aws_lambda_permission" "allow_cw_event" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "events.amazonaws.com"
}

resource "aws_cloudwatch_event_rule" "every_5_minutes" {
  name                = "${var.project}-every-5-minutes"
  schedule_expression = "rate(5 minutes)"
}

resource "aws_cloudwatch_event_target" "invoke_lambda" {
  rule      = aws_cloudwatch_event_rule.every_5_minutes.name
  target_id = "InvokeLambdaFunction"
  arn       = aws_lambda_function.this.arn
}

resource "aws_cloudwatch_log_group" "lambda_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.this.function_name}"
  retention_in_days = 14
}