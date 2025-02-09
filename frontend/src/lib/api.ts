import type { FileInfo } from './types';

const API_URL = import.meta.env.VITE_API_URL || 'http://40.90.193.108:8080/api';

function getAuthHeader(): HeadersInit {
    const token = localStorage.getItem('authToken');
    return token ? { 'Authorization': `Bearer ${token}` } : {};
}

export async function login(username: string, password: string) {
    const response = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
    });
    return response;
}

export async function signup(username: string, password: string, inviteCode: string) {
    const response = await fetch(`${API_URL}/signup`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password, invite: inviteCode })
    });
    return response;
}

export async function verifyInvite(code: string) {
    const response = await fetch(`${API_URL}/verify-invite`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ code })
    });
    return response;
}

export async function uploadFile(file: File): Promise<FileInfo> {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await fetch(`${API_URL}/upload`, {
        method: 'POST',
        headers: getAuthHeader(),
        body: formData
    });
    
    if (!response.ok) {
        throw new Error('Upload failed');
    }
    
    return response.json();
}

export async function listFiles(): Promise<FileInfo[]> {
    const response = await fetch(`${API_URL}/files`, {
        headers: getAuthHeader()
    });
    
    if (!response.ok) {
        throw new Error('Failed to fetch files');
    }
    
    return response.json();
}