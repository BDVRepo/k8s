#!/bin/bash

# Function to print table header
print_header() {
    printf "%-20s %-40s %-25s\n" "Namespace" "Resource/Name" "Created At"
    printf "%-20s %-40s %-25s\n" "---------" "-------------" "----------"
}

# Function to print a row in the table
print_row() {
    local namespace=$1
    local resource_name=$2
    local created_at=$3
    printf "%-20s %-40s %-25s\n" "$namespace" "$resource_name" "$created_at"
}

# Print table header
print_header

# Fetch NodePort and LoadBalancer services
kubectl get svc --all-namespaces -o jsonpath='{range .items[*]}{.metadata.namespace}{"\t"}{.metadata.name}{"\t"}{.spec.type}{"\t"}{.metadata.creationTimestamp}{"\n"}{end}' | grep -E "NodePort|LoadBalancer" | while read -r line; do
    namespace=$(echo "$line" | awk '{print $1}')
    name=$(echo "$line" | awk '{print $2}')
    resource_type=$(echo "$line" | awk '{print $3}')
    created_at=$(echo "$line" | awk '{print $4}')
    print_row "$namespace" "${resource_type}/${name}" "$created_at"
done

# Fetch Ingresses
kubectl get ing --all-namespaces -o jsonpath='{range .items[*]}{.metadata.namespace}{"\t"}{.metadata.name}{"\t"}{.metadata.creationTimestamp}{"\n"}{end}' | while read -r line; do
    namespace=$(echo "$line" | awk '{print $1}')
    name=$(echo "$line" | awk '{print $2}')
    created_at=$(echo "$line" | awk '{print $3}')
    print_row "$namespace" "ing/${name}" "$created_at"
done

# Fetch Pods with host networking
kubectl get pods --all-namespaces -o jsonpath='{range .items[*]}{.metadata.namespace}{"\t"}{.metadata.name}{"\t"}{.spec.hostNetwork}{"\t"}{.metadata.creationTimestamp}{"\n"}{end}' | grep "true" | while read -r line; do
    namespace=$(echo "$line" | awk '{print $1}')
    name=$(echo "$line" | awk '{print $2}')
    created_at=$(echo "$line" | awk '{print $4}')
    print_row "$namespace" "pod/${name}" "$created_at"
done
