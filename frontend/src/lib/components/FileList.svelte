<script lang="ts">
    import { onMount } from 'svelte';
    import { listFiles } from '$lib/api';
    import type { FileInfo } from '$lib/types';
    
    let files: FileInfo[] = [];
    let loading = true;
    let error: string | null = null;
    
    // Export this function so it can be called from outside
    export async function loadFiles() {
        loading = true;
        try {
            files = await listFiles();
        } catch (e) {
            error = 'Failed to load files';
        } finally {
            loading = false;
        }
    }
    
    onMount(loadFiles); // Call it on mount too
</script>

{#if loading}
    <div>Loading...</div>
{:else if error}
    <div class="error">{error}</div>
{:else}
    <div class="files-grid">
        {#each files as file}
            <div class="file-card">
                <h3>{file.name}</h3>
                <p>{new Date(file.uploadDate).toLocaleString()}</p>
                <p>{(file.size / 1024 / 1024).toFixed(2)} MB</p>
                <a href={file.url} target="_blank" rel="noopener noreferrer">
                    Download
                </a>
            </div>
        {/each}
    </div>
{/if}

<style>
    .error {
        color: #e53e3e;
        padding: 1rem;
        background: #fed7d7;
        border-radius: 4px;
        margin: 1rem 0;
    }
</style>