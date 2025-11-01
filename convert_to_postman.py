#!/usr/bin/env python3
"""
HTTP to Postman Collection Converter
Converts .http files to Postman Collection v2.1 format
"""

import json
import os
import re
from pathlib import Path
from typing import Dict, List, Any


def parse_http_file(file_path: str) -> List[Dict[str, Any]]:
    """Parse a .http file and extract all requests."""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # Split by ### separator
    requests = []
    sections = re.split(r'\n###\s*\n', content)

    for section in sections:
        if not section.strip():
            continue

        lines = section.strip().split('\n')
        if len(lines) < 2:
            continue

        # First line is the request name
        name = lines[0].replace('###', '').strip()
        if not name:
            continue

        # Find the HTTP method line
        method_line_idx = None
        for i, line in enumerate(lines[1:], 1):
            if re.match(r'^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\s+', line):
                method_line_idx = i
                break

        if method_line_idx is None:
            continue

        # Parse method and URL
        method_line = lines[method_line_idx]
        parts = method_line.split(None, 1)
        if len(parts) < 2:
            continue

        method = parts[0]
        url = parts[1]

        # Parse headers
        headers = []
        body_start_idx = None

        for i in range(method_line_idx + 1, len(lines)):
            line = lines[i].strip()

            # Empty line indicates start of body
            if not line:
                body_start_idx = i + 1
                break

            # Parse header
            if ':' in line:
                key, value = line.split(':', 1)
                headers.append({
                    'key': key.strip(),
                    'value': value.strip()
                })

        # Parse body
        body = None
        if body_start_idx and body_start_idx < len(lines):
            body_lines = lines[body_start_idx:]
            body_text = '\n'.join(body_lines).strip()
            if body_text and body_text.startswith('{'):
                body = body_text

        requests.append({
            'name': name,
            'method': method,
            'url': url,
            'headers': headers,
            'body': body
        })

    return requests


def create_postman_request(req: Dict[str, Any]) -> Dict[str, Any]:
    """Convert parsed request to Postman request format."""
    postman_req = {
        'name': req['name'],
        'request': {
            'method': req['method'],
            'header': req['headers'],
            'url': {
                'raw': req['url'],
                'protocol': 'http',
                'host': ['localhost'],
                'port': '8080',
                'path': []
            }
        },
        'response': []
    }

    # Parse URL
    url = req['url']
    if '://' in url:
        parts = url.split('://', 1)[1]
        if '/' in parts:
            host_port, path = parts.split('/', 1)
            path_parts = path.split('?')[0].split('/')
            postman_req['request']['url']['path'] = path_parts

            # Add query parameters
            if '?' in url:
                query_string = url.split('?', 1)[1]
                query_params = []
                for param in query_string.split('&'):
                    if '=' in param:
                        key, value = param.split('=', 1)
                        query_params.append({
                            'key': key,
                            'value': value
                        })
                postman_req['request']['url']['query'] = query_params

    # Add body if present
    if req['body']:
        postman_req['request']['body'] = {
            'mode': 'raw',
            'raw': req['body'],
            'options': {
                'raw': {
                    'language': 'json'
                }
            }
        }

    # Add test script for login/register endpoints to capture tokens
    if 'login' in req['name'].lower() or 'register' in req['name'].lower():
        if req['method'] == 'POST':
            postman_req['event'] = [{
                'listen': 'test',
                'script': {
                    'exec': [
                        'if (pm.response.code === 200 || pm.response.code === 201) {',
                        '    const response = pm.response.json();',
                        '    if (response.data && response.data.access_token) {',
                        '        pm.environment.set("access_token", response.data.access_token);',
                        '    }',
                        '    if (response.data && response.data.refresh_token) {',
                        '        pm.environment.set("refresh_token", response.data.refresh_token);',
                        '    }',
                        '}'
                    ],
                    'type': 'text/javascript'
                }
            }]

    return postman_req


def convert_http_to_postman(api_test_dir: str) -> Dict[str, Any]:
    """Convert all .http files in directory to Postman collection."""

    collection = {
        'info': {
            'name': 'DFood API',
            'description': 'Complete API collection for DFood backend',
            'schema': 'https://schema.getpostman.com/json/collection/v2.1.0/collection.json'
        },
        'item': []
    }

    # Get all .http files recursively
    http_files = sorted(Path(api_test_dir).glob('**/*.http'))

    # Group files by folder
    folders = {}
    for http_file in http_files:
        # Get the parent folder name relative to api_test_dir
        relative_path = http_file.relative_to(api_test_dir)
        folder_name = relative_path.parent.name if relative_path.parent.name else http_file.stem

        # Parse file
        requests = parse_http_file(str(http_file))
        if not requests:
            continue

        # Initialize folder if not exists
        if folder_name not in folders:
            folders[folder_name] = []

        # Add requests to folder
        folders[folder_name].extend([create_postman_request(req) for req in requests])

    # Create collection structure
    for folder_name, requests in sorted(folders.items()):
        folder = {
            'name': folder_name.replace('-', ' ').title(),
            'item': requests
        }
        collection['item'].append(folder)

    return collection


def main():
    """Main function."""
    script_dir = Path(__file__).parent
    api_test_dir = script_dir / 'api-test'
    output_file = script_dir / 'postman-collection.json'

    print(f"Converting .http files from: {api_test_dir}")

    if not api_test_dir.exists():
        print(f"Error: Directory not found: {api_test_dir}")
        return

    # Convert files
    collection = convert_http_to_postman(str(api_test_dir))

    # Write output
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(collection, f, indent=2)

    print(f"\n✅ Success!")
    print(f"Created: {output_file}")
    print(f"Total folders: {len(collection['item'])}")

    total_requests = sum(len(folder['item']) for folder in collection['item'])
    print(f"Total requests: {total_requests}")

    print("\n📋 Next steps:")
    print("1. Open Postman")
    print("2. Click 'Import' button")
    print("3. Select 'postman-collection.json'")
    print("4. Create environment 'DFood Local' with:")
    print("   - base_url: http://localhost:8080/api/v1")
    print("   - access_token: (will auto-fill on login)")
    print("   - refresh_token: (will auto-fill on login)")


if __name__ == '__main__':
    main()
