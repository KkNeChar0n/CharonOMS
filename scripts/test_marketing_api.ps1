# Marketing Management Module API Testing Script (PowerShell)
# This script tests all endpoints for activity templates and activities

# Configuration
$BaseUrl = "http://localhost:8080/api"
$Token = ""
$TemplateId = 0
$ActivityId = 0

# Function to print test results
function Print-Test {
    param([string]$Message)
    Write-Host "`n========================================" -ForegroundColor Yellow
    Write-Host "TEST: $Message" -ForegroundColor Yellow
    Write-Host "========================================" -ForegroundColor Yellow
}

function Print-Success {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Print-Error {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor Red
}

# Step 1: Login to get JWT token
Print-Test "Login to get JWT token"
try {
    $LoginBody = @{
        username = "admin"
        password = "admin123"
    } | ConvertTo-Json

    $LoginResponse = Invoke-RestMethod -Uri "$BaseUrl/login" -Method Post -Body $LoginBody -ContentType "application/json"

    $Token = $LoginResponse.data.token

    if ($Token) {
        Print-Success "Login successful, token obtained"
        Write-Host "Token: $($Token.Substring(0, [Math]::Min(50, $Token.Length)))..."
    } else {
        Print-Error "Login failed - No token received"
        exit 1
    }
} catch {
    Print-Error "Login failed: $_"
    exit 1
}

$Headers = @{
    "Authorization" = "Bearer $Token"
    "Content-Type" = "application/json"
}

# Step 2: Create Activity Template (按分类)
Print-Test "Create Activity Template (按分类)"
try {
    $CreateTemplateBody = @{
        name = "测试满减活动模板"
        type = 1
        select_type = 1
        classify_ids = @(1, 2)
        status = 0
    } | ConvertTo-Json

    $CreateTemplateResponse = Invoke-RestMethod -Uri "$BaseUrl/activity-templates" -Method Post -Body $CreateTemplateBody -Headers $Headers

    $TemplateId = $CreateTemplateResponse.data.id

    if ($TemplateId) {
        Print-Success "Template created successfully, ID: $TemplateId"
    } else {
        Print-Error "Template creation failed"
    }
} catch {
    Print-Error "Template creation failed: $_"
}

# Step 3: Get Activity Templates List
Print-Test "Get Activity Templates List"
try {
    $TemplatesList = Invoke-RestMethod -Uri "$BaseUrl/activity-templates" -Method Get -Headers $Headers

    $TemplatesCount = $TemplatesList.data.Count
    Print-Success "Retrieved $TemplatesCount templates"
    $TemplatesList.data | ConvertTo-Json -Depth 10
} catch {
    Print-Error "Failed to get templates list: $_"
}

# Step 4: Get Active Templates
Print-Test "Get Active Templates"
try {
    $ActiveTemplates = Invoke-RestMethod -Uri "$BaseUrl/activity-templates/active" -Method Get -Headers $Headers

    $ActiveCount = $ActiveTemplates.data.Count
    Print-Success "Retrieved $ActiveCount active templates"
    $ActiveTemplates.data | ConvertTo-Json -Depth 10
} catch {
    Print-Error "Failed to get active templates: $_"
}

# Step 5: Get Template Detail
Print-Test "Get Template Detail (ID: $TemplateId)"
try {
    $TemplateDetail = Invoke-RestMethod -Uri "$BaseUrl/activity-templates/$TemplateId" -Method Get -Headers $Headers

    Print-Success "Template detail retrieved"
    $TemplateDetail.data | ConvertTo-Json -Depth 10
} catch {
    Print-Error "Failed to get template detail: $_"
}

# Step 6: Update Template
Print-Test "Update Template (ID: $TemplateId)"
try {
    $UpdateTemplateBody = @{
        name = "更新后的满减活动模板"
        type = 1
        select_type = 1
        classify_ids = @(1, 2, 3)
        status = 0
    } | ConvertTo-Json

    $UpdateTemplateResponse = Invoke-RestMethod -Uri "$BaseUrl/activity-templates/$TemplateId" -Method Put -Body $UpdateTemplateBody -Headers $Headers

    if ($UpdateTemplateResponse.code -eq 0) {
        Print-Success "Template updated successfully"
    } else {
        Print-Error "Template update failed"
    }
} catch {
    Print-Error "Template update failed: $_"
}

# Step 7: Create Activity Template (按商品)
Print-Test "Create Activity Template (按商品)"
try {
    $CreateTemplate2Body = @{
        name = "测试满折活动模板"
        type = 2
        select_type = 2
        goods_ids = @(1, 2, 3)
        status = 0
    } | ConvertTo-Json

    $CreateTemplate2Response = Invoke-RestMethod -Uri "$BaseUrl/activity-templates" -Method Post -Body $CreateTemplate2Body -Headers $Headers

    $TemplateId2 = $CreateTemplate2Response.data.id

    if ($TemplateId2) {
        Print-Success "Template 2 created successfully, ID: $TemplateId2"
    } else {
        Print-Error "Template 2 creation failed"
    }
} catch {
    Print-Error "Template 2 creation failed: $_"
}

# Step 8: Create Activity
Print-Test "Create Activity"
try {
    $CreateActivityBody = @{
        name = "测试满减活动"
        template_id = $TemplateId
        start_time = "2026-02-01T00:00:00Z"
        end_time = "2026-02-28T23:59:59Z"
        status = 0
        details = @(
            @{
                threshold_amount = 100.0
                discount_value = 10.0
            },
            @{
                threshold_amount = 200.0
                discount_value = 25.0
            }
        )
    } | ConvertTo-Json -Depth 10

    $CreateActivityResponse = Invoke-RestMethod -Uri "$BaseUrl/activities" -Method Post -Body $CreateActivityBody -Headers $Headers

    $ActivityId = $CreateActivityResponse.data.id

    if ($ActivityId) {
        Print-Success "Activity created successfully, ID: $ActivityId"
    } else {
        Print-Error "Activity creation failed"
    }
} catch {
    Print-Error "Activity creation failed: $_"
}

# Step 9: Get Activities List
Print-Test "Get Activities List"
try {
    $ActivitiesList = Invoke-RestMethod -Uri "$BaseUrl/activities" -Method Get -Headers $Headers

    $ActivitiesCount = $ActivitiesList.data.Count
    Print-Success "Retrieved $ActivitiesCount activities"
    $ActivitiesList.data | ConvertTo-Json -Depth 10
} catch {
    Print-Error "Failed to get activities list: $_"
}

# Step 10: Get Activity Detail
Print-Test "Get Activity Detail (ID: $ActivityId)"
try {
    $ActivityDetail = Invoke-RestMethod -Uri "$BaseUrl/activities/$ActivityId" -Method Get -Headers $Headers

    Print-Success "Activity detail retrieved"
    $ActivityDetail.data | ConvertTo-Json -Depth 10
} catch {
    Print-Error "Failed to get activity detail: $_"
}

# Step 11: Get Activities by Date Range
Print-Test "Get Activities by Date Range"
try {
    $ActivitiesByDate = Invoke-RestMethod -Uri "$BaseUrl/activities/by-date-range?payment_time=2026-02-15T12:00:00Z" -Method Get -Headers $Headers

    Print-Success "Activities by date range retrieved"
    $ActivitiesByDate.data | ConvertTo-Json -Depth 10
} catch {
    Print-Error "Failed to get activities by date range: $_"
}

# Step 12: Update Activity
Print-Test "Update Activity (ID: $ActivityId)"
try {
    $UpdateActivityBody = @{
        name = "更新后的满减活动"
        template_id = $TemplateId
        start_time = "2026-02-01T00:00:00Z"
        end_time = "2026-02-28T23:59:59Z"
        status = 0
        details = @(
            @{
                threshold_amount = 150.0
                discount_value = 20.0
            },
            @{
                threshold_amount = 300.0
                discount_value = 50.0
            }
        )
    } | ConvertTo-Json -Depth 10

    $UpdateActivityResponse = Invoke-RestMethod -Uri "$BaseUrl/activities/$ActivityId" -Method Put -Body $UpdateActivityBody -Headers $Headers

    if ($UpdateActivityResponse.code -eq 0) {
        Print-Success "Activity updated successfully"
    } else {
        Print-Error "Activity update failed"
    }
} catch {
    Print-Error "Activity update failed: $_"
}

# Step 13: Update Activity Status
Print-Test "Update Activity Status (ID: $ActivityId)"
try {
    $UpdateStatusBody = @{
        status = 1
    } | ConvertTo-Json

    $UpdateStatusResponse = Invoke-RestMethod -Uri "$BaseUrl/activities/$ActivityId/status" -Method Put -Body $UpdateStatusBody -Headers $Headers

    if ($UpdateStatusResponse.code -eq 0) {
        Print-Success "Activity status updated successfully"
    } else {
        Print-Error "Activity status update failed"
    }
} catch {
    Print-Error "Activity status update failed: $_"
}

# Step 14: Delete Activity
Print-Test "Delete Activity (ID: $ActivityId)"
try {
    $DeleteActivityResponse = Invoke-RestMethod -Uri "$BaseUrl/activities/$ActivityId" -Method Delete -Headers $Headers

    if ($DeleteActivityResponse.code -eq 0) {
        Print-Success "Activity deleted successfully"
    } else {
        Print-Error "Activity deletion failed"
    }
} catch {
    Print-Error "Activity deletion failed: $_"
}

# Step 15: Update Template Status
Print-Test "Update Template Status (ID: $TemplateId)"
try {
    $UpdateTemplateStatusBody = @{
        status = 1
    } | ConvertTo-Json

    $UpdateTemplateStatusResponse = Invoke-RestMethod -Uri "$BaseUrl/activity-templates/$TemplateId/status" -Method Put -Body $UpdateTemplateStatusBody -Headers $Headers

    if ($UpdateTemplateStatusResponse.code -eq 0) {
        Print-Success "Template status updated successfully"
    } else {
        Print-Error "Template status update failed"
    }
} catch {
    Print-Error "Template status update failed: $_"
}

# Step 16: Delete Template (should fail if there are related activities)
Print-Test "Delete Template (ID: $TemplateId)"
try {
    $DeleteTemplateResponse = Invoke-RestMethod -Uri "$BaseUrl/activity-templates/$TemplateId" -Method Delete -Headers $Headers

    if ($DeleteTemplateResponse.code -eq 0) {
        Print-Success "Template deleted successfully"
    } else {
        Print-Success "Template deletion prevented (expected if has activities)"
    }
} catch {
    Print-Success "Template deletion prevented (expected if has activities)"
    Write-Host $_.Exception.Message
}

# Step 17: Delete Template 2 (if it was created)
if ($TemplateId2) {
    Print-Test "Delete Template 2 (ID: $TemplateId2)"
    try {
        $DeleteTemplate2Response = Invoke-RestMethod -Uri "$BaseUrl/activity-templates/$TemplateId2" -Method Delete -Headers $Headers

        if ($DeleteTemplate2Response.code -eq 0) {
            Print-Success "Template 2 deleted successfully"
        } else {
            Print-Error "Template 2 deletion failed"
        }
    } catch {
        Print-Error "Template 2 deletion failed: $_"
    }
}

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "ALL TESTS COMPLETED" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
