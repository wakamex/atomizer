"""Utility functions for loading environment variables."""

import os


def load_env_vars(env_path=".env", target_keys=None):
    """Load environment variables from a .env file.
    
    Args:
        env_path: Path to the .env file
        target_keys: List of keys to load. If None, loads all keys.
        
    Returns:
        List of loaded variable names
    """
    loaded_vars = []
    try:
        with open(env_path, "r", encoding="utf-8") as f:
            for line in f:
                stripped_line = line.strip()
                if not stripped_line or stripped_line.startswith("#") or "=" not in stripped_line:
                    continue
                key, value = stripped_line.split("=", 1)
                key = key.strip()
                value = value.strip().strip("'\"")  # Remove potential quotes
                
                # If target_keys is provided, only load those keys
                if target_keys is None or key in target_keys:
                    os.environ[key] = value
                    if key not in loaded_vars:
                        loaded_vars.append(key)
        
        if loaded_vars:
            print(f"Loaded {', '.join(loaded_vars)} from {env_path}")
    except FileNotFoundError:
        print(f"Warning: {env_path} file not found. Required variables should be set directly in your environment.")
    except Exception as e:
        print(f"Warning: Error reading {env_path}: {e}")
    
    return loaded_vars
