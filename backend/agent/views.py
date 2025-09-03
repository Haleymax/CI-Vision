from django.http import JsonResponse
from django.views.decorators.http import require_http_methods, require_GET
from django.views.decorators.csrf import csrf_exempt
from datetime import datetime
import json
import os

from utils.ai_agent import FileAgent

# Create your views here.

@require_GET
@csrf_exempt
def hello_world(request):
    """
    返回Hello World的GET接口
    
    Returns:
        JsonResponse: 包含问候信息的JSON响应
    """
    return JsonResponse({
        'message': 'Hello World',
        'status': 'success',
        'timestamp': datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
        'method': request.method,
        'path': request.path,
        'user_agent': request.META.get('HTTP_USER_AGENT', 'Unknown')
    }, json_dumps_params={'ensure_ascii': False, 'indent': 2})


@require_GET
@csrf_exempt
def list_files_with_ai(request):
    """
    使用file_tools和FileAgent列出文件
    """
    try:
        base_path = os.path.dirname(os.path.dirname(__file__))

        base_path = os.path.dirname(os.path.dirname(__file__))
        api_key = os.getenv('DEEPSEEK_API_KEY', 'sk-dea088bef52f4ef4a1d0ba4feeaf4ed6')
        file_agent = FileAgent(model="deepseek-chat", api_key=api_key)
        result = file_agent.list_all_files(base_path)
        
        return JsonResponse({
            'success': True,
            'base_path': base_path,
            "result": result
        }, json_dumps_params={'ensure_ascii': False, 'indent': 2})
        
    except Exception as e:
        return JsonResponse({
            'error': str(e),
            'timestamp': datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        }, status=500, json_dumps_params={'ensure_ascii': False, 'indent': 2})


@require_http_methods(["GET", "POST"])
@csrf_exempt
def api_demo(request):
    """
    演示GET和POST方法的API接口
    
    GET: 返回欢迎信息
    POST: 接收并回显数据
    """
    if request.method == 'GET':
        return JsonResponse({
            'message': '欢迎使用API演示接口',
            'methods_allowed': ['GET', 'POST'],
            'timestamp': datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        }, json_dumps_params={'ensure_ascii': False, 'indent': 2})
    
    elif request.method == 'POST':
        try:
            # 尝试解析JSON数据
            if request.content_type == 'application/json':
                data = json.loads(request.body)
            else:
                data = dict(request.POST)
            
            return JsonResponse({
                'message': '数据接收成功',
                'received_data': data,
                'timestamp': datetime.now().strftime('%Y-%m-%d %H:%M:%S')
            }, json_dumps_params={'ensure_ascii': False, 'indent': 2})
            
        except json.JSONDecodeError:
            return JsonResponse({
                'error': 'Invalid JSON data'
            }, status=400, json_dumps_params={'ensure_ascii': False, 'indent': 2})
