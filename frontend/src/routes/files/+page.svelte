<script lang="ts">
    import FileUpload from '$lib/components/FileUpload.svelte';
    import FileList from '$lib/components/FileList.svelte';
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    
    let fileListComponent: FileList;

    onMount(() => {
        // Check if user is authenticated
        const token = localStorage.getItem('authToken');
        if (!token) {
            goto('/login');
        }
    });
</script>

<main>
    <div class="header">
        <h1>File Storage</h1>
        <button 
            class="logout" 
            on:click={() => {
                localStorage.removeItem('authToken');
                goto('/login');
            }}
        >
            Logout
        </button>
    </div>
    <FileUpload onUploadComplete={() => fileListComponent.loadFiles()} />
    <FileList bind:this={fileListComponent} />
</main>

<style>
    main {
        width: 95%;
        max-width: 1200px;
        margin: 0 auto;
        padding: 1rem;
    }
    
    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 2rem;
    }

    h1 {
        margin: 0;
        font-size: clamp(1.5rem, 4vw, 2.5rem);
    }

    .logout {
        padding: 0.5rem 1rem;
        background: #dc3545;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }

    .logout:hover {
        background: #c82333;
    }
</style>