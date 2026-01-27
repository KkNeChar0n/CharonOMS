@echo off
chcp 65001 >nul
echo ========================================
echo 密码迁移脚本
echo ========================================
echo.
echo 此脚本将把数据库中的明文密码转换为bcrypt加密
echo.
echo 请确保已设置以下环境变量（或使用默认值）：
echo   DB_HOST (默认: localhost)
echo   DB_PORT (默认: 3306)
echo   DB_USER (默认: root)
echo   DB_PASSWORD (默认: 空)
echo   DB_NAME (默认: charonoms)
echo.
pause
echo.
echo 开始迁移...
echo.

cd /d "%~dp0.."
go run scripts/migrate_passwords.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo ✅ 密码迁移成功！
    echo ========================================
) else (
    echo.
    echo ========================================
    echo ❌ 密码迁移失败，请检查错误信息
    echo ========================================
)

echo.
pause
