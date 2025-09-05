# Удалите текущий файл и создайте заново с правильным форматом
rm docker/run_command.sh

# Создайте новый файл с Unix форматом строк
echo '#!/bin/bash

safe_exec() {
    local cmd="$1"
    local timeout=${2:-30}
    
    local forbidden_commands=(
        "rm -rf /"
        "rm -rf /*"
        ":(){:|:&};:"
        "mkfs"
        "dd if=/dev/zero"
        "chmod -R 777 /"
        "shutdown"
        "reboot"
        "halt"
        "poweroff"
    )
    
    for forbidden in "${forbidden_commands[@]}"; do
        if [[ "$cmd" == *"$forbidden"* ]]; then
            echo "ERROR: Command not allowed for security reasons"
            exit 1
        fi
    done
    
    timeout $timeout bash -c "$cmd" 2>&1
    local exit_code=$?
    
    if [ $exit_code -eq 124 ]; then
        echo "ERROR: Command timed out after $timeout seconds"
    elif [ $exit_code -ne 0 ]; then
        echo "ERROR: Command failed with exit code $exit_code"
    fi
    
    return $exit_code
}

if [ $# -eq 0 ]; then
    echo "Usage: $0 <command> [timeout]"
    exit 1
fi

command="$1"
timeout=${2:-30}

safe_exec "$command" $timeout' > docker/run_command.sh