<script setup lang="ts">
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import { AddTorrentFromFileContent, AddTorrentFromString } from "../../wailsjs/go/main/App";
import { arrayBufferToArrayNumber } from "@/ultis";

const router = useRouter();

interface HTMLInputEvent extends Event {
  target: HTMLInputElement & EventTarget;
}

const textInput = ref("");
const fileInput = ref<File | null>(null);
const fileInputRef = ref<HTMLInputElement | null>(null);

const onTextInputChange = (e: Event) => {
  textInput.value = (e.target as HTMLInputElement).value;
  fileInput.value = null;
  fileInputRef.value!.value = ""; // reset file input value
};

const onFileInputChange = (ev: Event) => {
  let files = (ev as HTMLInputEvent).target.files || (ev as DragEvent).dataTransfer?.files;
  if (!files?.length) {
    fileInput.value = null;
    return;
  }

  textInput.value = "";
  fileInput.value = files[0];
};

const getInfo = async () => {
  let infoHash = "";
  if (textInput.value) {
    console.log("text");
    infoHash = await AddTorrentFromString(textInput.value);
  } else if (fileInput.value) {
    console.log("file");
    const arrBuf = await fileInput.value.arrayBuffer();
    infoHash = await AddTorrentFromFileContent(arrayBufferToArrayNumber(arrBuf));
  } else {
    return;
  }

  console.log(infoHash);
  await router.push({ name: "files", params: { infoHash: infoHash } });
};

const disableGetInfo = computed<boolean>(() => {
  return !textInput.value && !fileInput.value;
});

</script>

<template>
  <main class="flex flex-col justify-center items-center w-full flex-grow h-dvh">
    <img src="../assets/logo.png" alt="TorPlayer logo" class="object-scale-down h-48 w-48"/>
    <h1 class="text-4xl font-bold mb-2">TorPlayer</h1>
    <p>A simple torrent player built with Go and Wails!</p>

    <div class="flex flex-col justify-between items-center mt-4">
      <label for="textInput" class="block text-sm mb-2">Input Torrent source</label>
      <input
        type="text"
        id="textInput"
        name="textInput"
        :value="textInput"
        @input="onTextInputChange"
        autocomplete="off"
        placeholder="Enter magnet, hash or torrent file link"
        class="p-4 rounded border border-gray-400 w-full bg-stone-700 text-stone-100 hover:border-red-700 hover:bg-stone-800"
      />

      <label for="fileInput" class="block text-sm my-2">Or Upload Torrent File</label>
      <input
        type="file"
        ref="fileInputRef"
        name="fileInput"
        accept=".torrent"
        @input="onFileInputChange"
        class="border text-stone-100 border-gray-400 hover:border-red-700 file:mr-5 rounded bg-stone-700 file:border-[0px] file:p-4 file:bg-stone-700 file:rounded-l file:text-stone-100 hover:file:cursor-pointer hover:file:bg-stone-800 hover:file:text-red-700"
      />
      <button class="px-6 py-2 my-3 bg-red-600 hover:bg-red-700 rounded disabled:bg-stone-500" @click="getInfo"
              :disabled="disableGetInfo">
        Get Info
      </button>
    </div>
  </main>
</template>

