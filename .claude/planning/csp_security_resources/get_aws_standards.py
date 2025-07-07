import boto3
import json
from botocore.exceptions import ClientError

def get_security_hub_standards():
    """
    Retrieve AWS Security Hub standards and format as JSON output
    with standard as an attribute for each entry.
    """
    try:
        # Initialize Security Hub client
        client = boto3.client('securityhub')
        
        # Get all standards
        response = client.describe_standards()
        
        # Process each standard and create JSON output
        standards_output = []
        
        for standard in response['Standards']:
            standard_json = {
                'standard': standard.get('StandardsArn', ''),
                'name': standard.get('Name', ''),
                'description': standard.get('Description', ''),
                'enabled_by_default': standard.get('EnabledByDefault', False),
                'standards_managed_by': standard.get('StandardsManagedBy', ''),
                'standard_id': standard.get('StandardsArn', '').split('/')[-1] if standard.get('StandardsArn') else ''
            }
            standards_output.append(standard_json)
        
        return standards_output
        
    except ClientError as e:
        print(f"Error retrieving Security Hub standards: {e}")
        return []
    except Exception as e:
        print(f"Unexpected error: {e}")
        return []

def main():
    """Main function to execute the script"""
    print("Retrieving AWS Security Hub standards...")
    
    standards = get_security_hub_standards()
    
    if standards:
        # Output each standard as a separate JSON object
        print("\n=== Security Hub Standards (Individual JSON Objects) ===")
        for standard in standards:
            print(json.dumps(standard, indent=2))
            print("-" * 50)
        
        # Also output as a single JSON array
        print("\n=== All Standards as JSON Array ===")
        print(json.dumps(standards, indent=2))
        
        # Save to file
        with open('securityhub_standards.json', 'w') as f:
            json.dump(standards, f, indent=2)
        print(f"\nOutput saved to securityhub_standards.json")
        
    else:
        print("No standards found or error occurred.")

if __name__ == "__main__":
    main()