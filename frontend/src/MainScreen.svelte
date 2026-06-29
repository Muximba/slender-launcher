<script lang="ts">
  import Select from "svelte-select";
  import logo from "./assets/images/logo-universal.png";
  import {
    ActiveDownload,
    DownloadMaps,
    DownloadPercent,
    DownloadedBytes,
    DownloadedFiles,
    LocalEnabled,
    NeedsUpdate,
    News,
    Play,
    Revision,
    ServerStatus,
    TotalBytes,
    TotalFiles,
    Update,
    Version,
  } from "../wailsjs/go/main/App.js";
  import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
  import { onMount } from "svelte";
  import PlayIcon from "./PlayIcon.svelte";
  import UpdateIcon from "./UpdateIcon.svelte";
  import DownloadIcon from "./DownloadIcon.svelte";
  import SettingsIcon from "./SettingsIcon.svelte";

  export let openSettings: () => void;

  let version: string = "";
  let revision: number = 0;
  let updating: boolean = false;
  let ready: boolean = false;
  let needsUpdate: boolean = false;

  let serverOnline: boolean = false;
  let statusLoaded: boolean = false;
  let playersOnline: number = 0;
  let news: { date: string; title: string; content: string; url: string }[] =
    [];

  function openNews(url: string) {
    if (url) BrowserOpenURL(url);
  }

  let progress: number = 0;
  let totalFiles: number = 0;
  let totalBytes: number = 0;
  let downloadedFiles: number = 0;
  let downloadedBytes: number = 0;
  let activeDownload: string = "";

  let mapKind = 0;

  let hasLocal = false;

  onMount(async () => {
    // status do servidor PRIMEIRO (rapido, um POST): senao o pill fica "offline"
    // durante toda a checagem de arquivos (que e lenta).
    ServerStatus().then((status) => {
      serverOnline = status.online;
      playersOnline = status.playersOnline;
      statusLoaded = true;
    });
    News().then((n) => (news = n));

    // atualiza o status do servidor a cada 30s
    setInterval(async () => {
      const s = await ServerStatus();
      serverOnline = s.online;
      playersOnline = s.playersOnline;
    }, 30000);

    // depois a checagem de versao/arquivos (mais demorada)
    revision = await Revision();
    version = await Version();
    needsUpdate = await NeedsUpdate();
    ready = true;
    hasLocal = await LocalEnabled();
  });

  function update() {
    totalFiles = 0;
    totalBytes = 0;
    downloadedBytes = 0;
    downloadedFiles = 0;
    void Update();
    updating = true;

    const interval = setInterval(async () => {
      totalFiles = await TotalFiles();
      totalBytes = await TotalBytes();
      downloadedBytes = await DownloadedBytes();
      downloadedFiles = await DownloadedFiles();
      activeDownload = await ActiveDownload();
      progress = await DownloadPercent();

      // so considera concluido depois que a lista de arquivos foi montada
      // (totalFiles > 0); senao 0 === 0 encerraria no 1o tick e a UI sumiria
      if (totalFiles > 0 && downloadedFiles >= totalFiles) {
        updating = false;
        needsUpdate = false;
        clearInterval(interval);
      }
    }, 1000);
  }

  function downloadMaps() {
    if (mapKind == null) return;
    totalFiles = 0;
    totalBytes = 0;
    downloadedBytes = 0;
    downloadedFiles = 0;
    void DownloadMaps(mapKind);
    updating = true;

    const interval = setInterval(async () => {
      totalFiles = await TotalFiles();
      totalBytes = await TotalBytes();
      downloadedBytes = await DownloadedBytes();
      downloadedFiles = await DownloadedFiles();
      activeDownload = await ActiveDownload();
      progress = await DownloadPercent();

      if (progress === 100) {
        updating = false;
        needsUpdate = false;
        clearInterval(interval);
      }
    }, 1000);
  }

  function formatBytes(bytes: number, decimals = 2) {
    if (!+bytes) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = [
      "Bytes",
      "KiB",
      "MiB",
      "GiB",
      "TiB",
      "PiB",
      "EiB",
      "ZiB",
      "YiB",
    ];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
  }

  function play() {
    ready = false;
    Play(false);
  }

  function playLocal() {
    ready = false;
    Play(true);
  }

  const mapTypes = [
    { value: 0, label: "Full w/ markers" },
    { value: 1, label: "Full w/o markers" },
    { value: 2, label: "Overlayed w/ markers" },
    { value: 3, label: "Overlayed w/o markers" },
    { value: 4, label: "Overlayed w/ markers (+PoI)" },
  ];
</script>

<button class="settings" on:click={openSettings} disabled={updating}>
  <SettingsIcon />
</button>
<div>
  <img alt="Logo" id="logo" src={logo} />
  <div class="server-status">
    <span class="dot" class:online={serverOnline}></span>
    {#if !statusLoaded}
      connecting...
    {:else if serverOnline}
      {playersOnline} players online
    {:else}
      Server offline
    {/if}
  </div>
  <div class="actions">
    <div>
      <h3>Play</h3>
      {#if updating}
        <button class="update" disabled>
          {#if totalFiles === 0}
            <div>Preparing...</div>
            <div>checking files</div>
          {:else}
            <div>{Math.floor(progress)}% &middot; {downloadedFiles}/{totalFiles}</div>
            <div>
              {formatBytes(downloadedBytes)} / {formatBytes(totalBytes)}
            </div>
          {/if}
        </button>
      {:else if !ready}
        <div class="checking">
          <div class="spinner"></div>
          <div>Verificando arquivos...</div>
        </div>
      {:else if needsUpdate}
        <div>
          <button class="update" on:click={update} disabled={!ready}>
            <UpdateIcon />
          </button>
          Update available.
        </div>
      {:else}
        <div>
          <div class="row">
            <button
              class="play"
              class:withLocal={hasLocal}
              disabled={!ready}
              on:click={play}
            >
              <PlayIcon />
              {#if hasLocal}
                Remote
              {:else}
                {version} + rev.{revision}
              {/if}
            </button>
            {#if hasLocal}
              <button
                class="play"
                class:local={hasLocal}
                disabled={!ready}
                on:click={playLocal}
              >
                <PlayIcon />
                Local
              </button>
            {/if}
          </div>
          {#if ready}Up to date.{:else}Loading...{/if}
        </div>
      {/if}
    </div>
    <!-- Muximba: download de mapa removido — o minimap.otmm ja vem embutido no client (otclient nao le os PNGs da tibiamaps). -->
  </div>

  {#if news.length}
    <div class="news">
      <div class="news-header">Novidades</div>
      {#each news.slice(0, 3) as item}
        <div
          class="news-item"
          class:clickable={item.url}
          role={item.url ? "link" : undefined}
          on:click={() => openNews(item.url)}
        >
          <div class="news-head">
            <strong>{item.title}</strong>
            <span class="news-date">{item.date}</span>
          </div>
          {#if item.content}
            <span class="news-content">{item.content}</span>
          {/if}
          {#if item.url}
            <span class="news-more">ver mais →</span>
          {/if}
        </div>
      {/each}
    </div>
  {/if}

  {#if updating}
    <div class="progress-section">
      <div class="progress-bar">
        <div class="progress" style="width: {progress}%" />
        <div class="active-download">{activeDownload}</div>
      </div>
    </div>
  {/if}
</div>

<style>
  .progress-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  div.progress-bar {
    position: relative;
    align-items: start;
    justify-content: start;
    width: 512px;
    height: 32px;
    background-color: #333333;
    border-radius: 8px;
    margin: 8px 0;
  }

  .active-download {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    color: white;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    padding: 0 8px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .progress {
    height: 100%;
    background-color: #016f4e;
    border-radius: 8px;
    transition: width 0.5s ease-in-out;
  }

  div {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }

  button {
    background: none;
    border: none;
    cursor: pointer;
    padding: 8px;
    width: 200px;
    height: 74px;
    color: white;
    border-radius: 8px;
    box-shadow: #333333 0px 0px 4px 0px;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
  }

  button.update {
    flex-direction: column;
    background-color: #f4b343;
  }
  button.update:disabled {
    flex-direction: column;
    padding: 12px;
  }

  button:disabled {
    color: #ccc;
    background-color: #333333;
    opacity: 0.75;
  }

  button.play {
    background-color: #016f4e;
  }

  #logo {
    display: block;
    width: 172px;
    height: auto;
    margin: auto;
    padding: 3% 0 0;
    object-fit: contain;
  }

  .actions {
    display: flex;
    flex-direction: row;
    align-items: start;
    gap: 8px;
    width: 100%;
  }

  .maps {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .map-select {
    width: 200px;
    --border-radius: 16px;
    --list-border-radius: 16px;
    --item-color: #4e3bf5;
    --item-hover-color: #4e3bf5;
    --placeholder-color: #4e3bf5;
    --selected-item-color: #4e3bf5;
  }

  h3 {
    margin: 0;
    padding: 0;
    font-size: 16px;
    text-shadow: 0 1px 4px rgba(0, 0, 0, 0.9), 0 0 2px rgba(0, 0, 0, 0.7);
  }
  /* texto solto (Up to date / Loading / Update available) sobre fundo claro */
  .actions > div > div {
    text-shadow: 0 1px 4px rgba(0, 0, 0, 0.9);
  }

  .maps button {
    width: 100%;
    height: 24px;
    background-color: #4e3bf5;
  }

  .withLocal {
    width: 90px;
    display: flex;
    flex-direction: column;
  }

  .play.local {
    background-color: #ba3bf5;
    width: 90px;
    display: flex;
    flex-direction: column;
  }

  .row {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    gap: 8px;
  }

  button.settings {
    position: absolute;
    top: 0;
    right: 0;
    width: 48px;
    height: 48px;
    margin: 8px;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    box-shadow: none;
  }

  .server-status {
    flex-direction: row;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    font-weight: 600;
    color: #f1ece1;
    background: rgba(12, 18, 30, 0.55);
    padding: 5px 14px;
    border-radius: 999px;
    margin: 6px 0 12px;
    box-shadow: 0 1px 5px rgba(0, 0, 0, 0.35);
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
  }
  .server-status .dot {
    width: 9px;
    height: 9px;
    border-radius: 50%;
    background-color: #d9534f;
  }
  .server-status .dot.online {
    background-color: #6fbf73;
    box-shadow: 0 0 6px #6fbf73;
  }

  .checking {
    flex-direction: column;
    gap: 8px;
    color: #e7ddcb;
    font-size: 13px;
    height: 74px;
    justify-content: center;
  }
  .spinner {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    border: 3px solid rgba(255, 255, 255, 0.18);
    border-top-color: #6fbf73;
    animation: mux-spin 0.8s linear infinite;
  }
  @keyframes mux-spin {
    to { transform: rotate(360deg); }
  }
  @media (prefers-reduced-motion: reduce) {
    .spinner { animation: none; }
  }

  .news {
    flex-direction: column;
    gap: 6px;
    width: 512px;
    margin: 12px 0;
  }
  .news-header {
    flex-direction: row;
    align-self: flex-start;
    font-size: 12px;
    font-weight: 700;
    letter-spacing: 1px;
    text-transform: uppercase;
    color: #e7ddcb;
    text-shadow: 0 1px 4px rgba(0, 0, 0, 0.9);
    margin-bottom: 2px;
  }
  .news-item {
    flex-direction: column;
    align-items: stretch;
    gap: 3px;
    width: 100%;
    background-color: rgba(42, 36, 54, 0.82);
    border-radius: 6px;
    padding: 8px 12px;
    color: #ece3d4;
    font-size: 13px;
    text-align: left;
  }
  .news-head {
    flex-direction: row;
    align-items: baseline;
    justify-content: space-between;
    width: 100%;
    gap: 8px;
  }
  .news-head strong {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .news-date {
    color: #9990a6;
    font-size: 12px;
    flex-shrink: 0;
  }
  .news-content {
    color: #b3a994;
    font-size: 12px;
    line-height: 1.35;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  .news-item.clickable {
    cursor: pointer;
    transition: background-color 0.15s ease;
  }
  .news-item.clickable:hover {
    background-color: rgba(62, 52, 82, 0.95);
  }
  .news-more {
    align-self: flex-end;
    font-size: 11px;
    font-weight: 600;
    color: #6fbf73;
  }
</style>
