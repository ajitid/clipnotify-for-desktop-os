WinRT/WinSDK requires a GUI window to be in foreground for clipboard monitoring to work (https://github.com/microsoft/WindowsAppSDK/issues/62). Thankfully we can use Win32 instead, which has an API for monitoring clipboard in CLI as well. Read more about them here:

I am using miniconda and I installed the package using `conda install pywin32`. Even though this project looks like it is using poetry, actually I'm not using it any manner — not even sourcing poetry's python environment. I am running the program using conda base like `python .\clipnotify_win\main.py`  
With poetry installing the package pywin32 happens successfully but the programs itself fails w/ error.

- https://stackoverflow.com/a/65857206/7683365
- https://abdus.dev/posts/monitor-clipboard/

  - https://github.com/abdusco/dumpclip
  - https://twitter.com/abdusdev

- https://www.ctrl.blog/entry/clipboard-security.html

  - http://nspasteboard.org/

- copy formats

  - https://superuser.com/questions/199285/how-to-copy-image-to-clipboard-to-paste-to-another-application
  - https://stackoverflow.com/questions/3571179/how-does-x11-clipboard-handle-multiple-data-formats
  - Make sure to educate your users when you implement any special clipboard handling.
    - tell users that telltale-sync automatically expires content after two minutes.
    - we copy stuff even if it is sensitive
      - you would think its not right (and that's true), but there's no universal way of transferring sensitive text and telling clipboard managers to not to store it
      - moreover, not copying everytime would mean that sometimes univ. clipb. works, sometimes it doesn't. And what about the case when you _actually_ want to receive sensitive data on your other machine? (that's why I think having a small 2min expiration is a good workaround).
    - copy image (local) -> receive text (telltail) -> restore << will give you nothing as we in the program we only copy and restore text
  - https://github.com/p0deje/Maccy#ignore-copied-items

- telltail sync should require a flag to allow autocopy for privacy reasons (autocopy requires some manual setup anyway). Copy should always be explicit

after copying an image from web browser, run this:
(needs focus on document, that's why i've used a settimeout so you can click and scroll on the page again)

```
setTimeout(async () => {
    const clipboardContents = await navigator.clipboard.read()

    for (const item of clipboardContents) {
      if (!item.types.includes("text/html")) {
        throw new Error("Clipboard contains non-text data.");
      }
      const blob = await item.getType("text/html");
      const text = await new Response(blob).text()
      console.log(text)
    }
}, 4000)
```

- xclip, xsel
  - xsel is only for plain text
  - xclip is also no go as transferring data back to xclip takes too much time because it doesn't support streaming https://www.reddit.com/r/linuxquestions/comments/x92ybc/comment/inn3x7f/?utm_source=share&utm_medium=web2x&context=3
  - all formats out
    - `xclip -selection clipboard -target TARGETS -out`
      - also has ways mentioned to get with GTK3 https://stackoverflow.com/questions/3571179/how-does-x11-clipboard-handle-multiple-data-formats
      - other stuff
        - xclip -selection CLIPBOARD -target image/png -in /path/to/some.png
  - errors
    - xclip on copying an image from chrome fails on `xclip -o -rmlastnl -selection clipboard`
      - also on `xclip -selection clipboard -target SAVE_TARGETS -out`
        - so because command can fail, so we should store as much as we can and skip the rest, and should also ensure removing that target from TARGETS
    - GTK on Xorg needs request mechanism https://stackoverflow.com/questions/3261379/getting-html-source-or-rich-text-from-the-x-clipboard/3263632#3263632
    - pbcopy ain't sufficient on mac to retrieve multiple formats, we need to use AppleScript https://superuser.com/a/1215618
- Windows via powershell https://stackoverflow.com/a/62377454
  - TL;DR https://www.pdq.com/powershell/get-clipboard/
  - https://learn.microsoft.com/en-us/powershell/module/Microsoft.PowerShell.Management/Get-Clipboard?view=powershell-5.1
  - NOTE: oddly, there isn't any MIME type listed, which is very, very odd. How will figma for example copy an artboard from one window to another if it can't specify its own mimetype?? Maybe true identification of all values isn't allowed (though -Raw can help, or maybe its' just only to give raw text string)
  - NOTE: sadly Get-Clipboard is even more barebones, not accepting anything other than text formats
