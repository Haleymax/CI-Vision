from pathlib import Path

def list_directory(path: str, recursive: bool = False) -> list:
    path = Path(path).resolve()
    
    if not path.exists():
        raise FileNotFoundError(f"Directory does not exist: {path}")
    
    if not path.is_dir():
        raise NotADirectoryError(f"Path is not a directory: {path}")
    
    files = []
    
    def _collect_files(current_path: Path):
        for item in current_path.iterdir():
            if item.is_file():
                files.append(str(item.absolute()))
            elif recursive and item.is_dir():
                _collect_files(item)
    
    _collect_files(path)
    return sorted(files)
