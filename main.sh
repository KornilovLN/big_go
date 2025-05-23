#!/bin/bash
# Script to manage big_go Docker environment

# Colors for better readability
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to display the header
show_header() {
    clear
    echo -e "${BLUE}=========================================${NC}"
    echo -e "${GREEN}   big_go Docker Management Script      ${NC}"
    echo -e "${BLUE}=========================================${NC}"
    echo ""
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo -e "${RED}Error: Docker is not running. Please start Docker and try again.${NC}"
        exit 1
    fi
}

# Function to check if docker-compose.yml exists
check_compose_file() {
    if [ ! -f "docker-compose.yml" ]; then
        echo -e "${RED}Error: docker-compose.yml not found in the current directory.${NC}"
        exit 1
    fi
}

# Function to create config files if they don't exist
create_config_files() {
    echo -e "${YELLOW}Checking configuration files...${NC}"

    if [ ! -f "config_rabbitmq.json" ]; then
        echo -e "${YELLOW}Creating config_rabbitmq.json...${NC}"
        cat > config_rabbitmq.json << EOF
{
  "host": "rabbitmq",
  "port": 5672,
  "user": "guest",
  "password": "guest",
  "vhost": "/"
}
EOF
    fi

    if [ ! -f "config_redis.json" ]; then
        echo -e "${YELLOW}Creating config_redis.json...${NC}"
        cat > config_redis.json << EOF
{
  "host": "redis",
  "port": 6379,
  "password": "",
  "db": 0
}
EOF
    fi

    if [ ! -f "config_postgresql.json" ]; then
        echo -e "${YELLOW}Creating config_postgresql.json...${NC}"
        cat > config_postgresql.json << EOF
{
  "host": "postgres",
  "port": 5432,
  "user": "postgres",
  "password": "postgres",
  "dbname": "big_go",
  "sslmode": "disable"
}
EOF
    fi

    if [ ! -f "config_go.json" ]; then
        echo -e "${YELLOW}Creating config_go.json...${NC}"
        cat > config_go.json << EOF
{
  "log_level": "info",
  "environment": "development"
}
EOF
    fi

    echo -e "${GREEN}All configuration files are ready.${NC}"
    sleep 1
}

# Function to start all services
start_all_services() {
    echo -e "${YELLOW}Starting all services...${NC}"
    create_config_files
    docker-compose up -d
    echo -e "${GREEN}All services started successfully!${NC}"
    echo -e "${BLUE}Веб-интерфейс RabbitMQ: http://localhost:15672 (логин: guest, пароль: guest)${NC}"
    echo -e "${BLUE}User1 Dashboard: http://localhost:8082${NC}"
    echo -e "${BLUE}User2 Dashboard: http://localhost:8083${NC}"
    sleep 2
}

# Function to stop all services
stop_all_services() {
    echo -e "${YELLOW}Stopping all services...${NC}"
    docker-compose down
    echo -e "${GREEN}All services stopped successfully!${NC}"
    sleep 2
}

# Function to restart all services
restart_all_services() {
    echo -e "${YELLOW}Restarting all services...${NC}"
    docker-compose restart
    echo -e "${GREEN}All services restarted successfully!${NC}"
    sleep 2
}

# Function to rebuild and restart all services
rebuild_all_services() {
    echo -e "${YELLOW}Rebuilding and restarting all services...${NC}"
    docker-compose down
    docker-compose build --no-cache
    docker-compose up -d
    echo -e "${GREEN}All services rebuilt and restarted successfully!${NC}"
    sleep 2
}

# Function to show service status
show_service_status() {
    echo -e "${YELLOW}Current service status:${NC}"
    docker-compose ps
    echo ""
    echo -e "${YELLOW}Press Enter to continue...${NC}"
    read
}

# Function to view logs
view_logs() {
    show_header
    echo -e "${YELLOW}Select a service to view logs (or 'all' for all services):${NC}"
    echo "1. All services"
    echo "2. generator"
    echo "3. collector"
    echo "4. user1"
    echo "5. user2"
    echo "6. postgres"
    echo "7. redis"
    echo "8. rabbitmq"
    echo "0. Back to main menu"

    read -p "Enter your choice: " log_choice

    case $log_choice in
        1|all|All|ALL)
            docker-compose logs
            ;;
        2|generator)
            docker-compose logs generator
            ;;
        3|collector)
            docker-compose logs collector
            ;;
        4|user1)
            docker-compose logs user1
            ;;
        5|user2)
            docker-compose logs user2
            ;;
        6|postgres)
            docker-compose logs postgres
            ;;
        7|redis)
            docker-compose logs redis
            ;;
        8|rabbitmq)
            docker-compose logs rabbitmq
            ;;
        0)
            return
            ;;
        *)
            echo -e "${RED}Invalid choice. Please try again.${NC}"
            sleep 1
            view_logs
            ;;
    esac

    echo -e "${YELLOW}Press Enter to continue...${NC}"
    read
}

# Function to manage specific services
manage_specific_services() {
    while true; do
        show_header
        echo -e "${YELLOW}Select a service to manage:${NC}"
        echo "1. generator"
        echo "2. collector"
        echo "3. user1"
        echo "4. user2"
        echo "5. postgres"
        echo "6. redis"
        echo "7. rabbitmq"
        echo "0. Back to main menu"

        read -p "Enter your choice: " service_choice

        case $service_choice in
            0)
                return
                ;;
            1|2|3|4|5|6|7)
                case $service_choice in
                    1) service="generator" ;;
                    2) service="collector" ;;
                    3) service="user1" ;;
                    4) service="user2" ;;
                    5) service="postgres" ;;
                    6) service="redis" ;;
                    7) service="rabbitmq" ;;
                esac

                show_header
                echo -e "${YELLOW}Managing service: ${GREEN}$service${NC}"
                echo "1. Start service"
                echo "2. Stop service"
                echo "3. Restart service"
                echo "4. Rebuild service"
                echo "5. View logs"
                echo "0. Back to service selection"

                read -p "Enter your choice: " action_choice

                case $action_choice in
                    1)
                        echo -e "${YELLOW}Starting $service...${NC}"
                        docker-compose up -d $service
                        echo -e "${GREEN}Service $service started successfully!${NC}"
                        sleep 2
                        ;;
                    2)
                        echo -e "${YELLOW}Stopping $service...${NC}"
                        docker-compose stop $service
                        echo -e "${GREEN}Service $service stopped successfully!${NC}"
                        sleep 2
                        ;;
                    3)
                        echo -e "${YELLOW}Restarting $service...${NC}"
                        docker-compose restart $service
                        echo -e "${GREEN}Service $service restarted successfully!${NC}"
                        sleep 2
                        ;;
                    4)
                        echo -e "${YELLOW}Rebuilding $service...${NC}"
                        docker-compose stop $service
                        docker-compose build --no-cache $service
                        docker-compose up -d $service
                        echo -e "${GREEN}Service $service rebuilt successfully!${NC}"
                        sleep 2
                        ;;
                    5)
                        docker-compose logs $service
                        echo -e "${YELLOW}Press Enter to continue...${NC}"
                        read
                        ;;
                    0)
                        ;;
                    *)
                        echo -e "${RED}Invalid choice. Please try again.${NC}"
                        sleep 1
                        ;;
                esac
                ;;
            *)
                echo -e "${RED}Invalid choice. Please try again.${NC}"
                sleep 1
                ;;
        esac
    done
}

# Function to clean up resources
cleanup_resources() {
    show_header
    echo -e "${YELLOW}Select cleanup option:${NC}"
    echo "1. Stop services and remove containers"
    echo "2. Stop services, remove containers and networks"
    echo "3. Stop services, remove containers, networks, and volumes (CAUTION: Data will be lost)"
    echo "4. Remove all unused containers, networks, and images (system prune)"
    echo "0. Back to main menu"

    read -p "Enter your choice: " cleanup_choice

    case $cleanup_choice in
        1)
            echo -e "${YELLOW}Stopping services and removing containers...${NC}"
            docker-compose down
            echo -e "${GREEN}Cleanup completed successfully!${NC}"
            sleep 2
            ;;
        2)
            echo -e "${YELLOW}Stopping services, removing containers and networks...${NC}"
            docker-compose down --remove-orphans
            echo -e "${GREEN}Cleanup completed successfully!${NC}"
            sleep 2
            ;;
        3)
            echo -e "${RED}WARNING: This will remove all data in volumes!${NC}"
            read -p "Are you sure you want to continue? (y/n): " confirm
            if [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]]; then
                echo -e "${YELLOW}Stopping services, removing containers, networks, and volumes...${NC}"
                docker-compose down -v --remove-orphans
                echo -e "${GREEN}Cleanup completed successfully!${NC}"
            else
                echo -e "${YELLOW}Operation cancelled.${NC}"
            fi
            sleep 2
            ;;
        4)
            echo -e "${RED}WARNING: This will remove all unused Docker resources!${NC}"
            read -p "Are you sure you want to continue? (y/n): " confirm
            if [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]]; then
                echo -e "${YELLOW}Removing all unused containers, networks, and images...${NC}"
                docker system prune -f
                echo -e "${GREEN}Cleanup completed successfully!${NC}"
            else
                echo -e "${YELLOW}Operation cancelled.${NC}"
            fi
            sleep 2
            ;;
        0)
            return
            ;;
        *)
            echo -e "${RED}Invalid choice. Please try again.${NC}"
            sleep 1
            cleanup_resources
            ;;
    esac
}

# Main function
main() {
    check_docker
    check_compose_file

    while true; do
        show_header
        echo -e "${YELLOW}Please select an option:${NC}"
        echo "1. Start all services"
        echo "2. Stop all services"
        echo "3. Restart all services"
        echo "4. Rebuild and restart all services"
        echo "5. Show service status"
        echo "6. View logs"
        echo "7. Manage specific services"
        echo "8. Cleanup resources"
        echo "0. Exit"

        read -p "Enter your choice: " choice

        case $choice in
            1)
                start_all_services
                ;;
            2)
                stop_all_services
                ;;
            3)
                restart_all_services
                ;;
            4)
                rebuild_all_services
                ;;
            5)
                show_service_status
                ;;
            6)
                view_logs
                ;;
            7)
                manage_specific_services
                ;;
            8)
                cleanup_resources
                ;;
            0)
                echo -e "${GREEN}Exiting. Goodbye!${NC}"
                exit 0
                ;;
            *)
                echo -e "${RED}Invalid choice. Please try again.${NC}"
                sleep 1
                ;;
        esac
    done
}

# Run the main function
main
