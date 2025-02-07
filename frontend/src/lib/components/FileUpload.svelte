<script lang="ts">
    import { uploadFile } from '$lib/api';
    
    export let onUploadComplete: () => void; 
    
    let files: FileList | null = null;
    let uploading = false;
    
    async function handleSubmit() {
        if (!files?.[0]) return;
        
        uploading = true;
        try {
            await uploadFile(files[0]);
            onUploadComplete();
        } catch (error) {
            console.error('Upload failed:', error);
        } finally {
            uploading = false;
            files = null;
        }
    }
</script>

<div class="upload-container">
    <input 
        type="file" 
        bind:files 
        disabled={uploading}
    />
    <button 
        on:click={handleSubmit} 
        disabled={!files?.[0] || uploading}
    >
        {uploading ? 'Uploading...' : 'Upload'}
    </button>
</div>

<style>
    .upload-container {
        margin: 2rem 0;
        padding: 2rem;
        border: 2px dashed #ccc;
        border-radius: 8px;
        text-align: center;
        background: #f9f9f9;
    }
    
    input[type="file"] {
        margin-right: 1rem;
    }
    
    button {
        padding: 0.75rem 1.5rem;
        background: #4a9eff;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        transition: background 0.2s;
    }
    
    button:hover:not(:disabled) {
        background: #2b7fd9;
    }
    
    button:disabled {
        background: #ccc;
        cursor: not-allowed;
    }
</style>