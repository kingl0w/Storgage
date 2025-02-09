<script lang="ts">
    import { onMount } from 'svelte';
    import { listFiles, deleteFile } from '$lib/api';
    import type { FileInfo } from '$lib/types';

    let files: FileInfo[] = [];
    let loading = true;
    let error: string | null = null;

    // Load files from server
    export async function loadFiles() {
        loading = true;
        error = null;
        try {
            files = await listFiles();
        } catch (e) {
            error = 'Failed to load files';
        } finally {
            loading = false;
        }
    }

    // Called when component mounts
    onMount(loadFiles);

    // Handle delete button click
    async function handleDelete(filename: string) {
        if (!confirm(`Are you sure you want to delete ${filename}?`)) {
            return;
        }

        try {
            const response = await deleteFile(filename);
            if (!response.ok) {
                // read the error message (could be JSON or plain text)
                let msg;
                try {
                    const data = await response.json();
                    msg = data.error || 'Delete failed';
                } catch (err) {
                    msg = await response.text();
                }
                error = msg;
            } else {
                // Remove the file from our array so the UI updates
                files = files.filter((f) => f.name !== filename);
            }
        } catch (err) {
            error = 'Delete request failed';
        }
    }
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
                <p>{(file.size / (1024 * 1024)).toFixed(2)} MB</p>

                <!-- Download Link -->
                <a
                    href={file.url}
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    Download
                </a>

                <!-- Delete Button -->
                <button on:click={() => handleDelete(file.name)}>
                    Delete
                </button>
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
    /* Example styles */
    .files-grid {
        display: flex;
        flex-wrap: wrap;
        gap: 1rem;
    }
    .file-card {
        border: 1px solid #ddd;
        border-radius: 4px;
        padding: 1rem;
        width: 200px;
    }
    button {
        margin-top: 0.5rem;
        cursor: pointer;
        background: #ef4444;
        color: #fff;
        border: none;
        padding: 0.5rem 0.75rem;
        border-radius: 4px;
    }
    button:hover {
        background: #dc2626;
    }
</style>
