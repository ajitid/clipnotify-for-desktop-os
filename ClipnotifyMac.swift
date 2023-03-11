// taken from https://github.com/p0deje/Maccy/blob/86b6dd79b2d0f52205bb044c2561175b683a6398/Maccy/Clipboard.swift#L61

import AppKit

class PasteboardWatcher {
    private let pasteboard = NSPasteboard.general
    
    private var changeCount : Int

    init() {
        changeCount = pasteboard.changeCount
    }

    func startListening () {
        Timer.scheduledTimer(timeInterval: 1.0, target: self, selector: #selector(self.checkForChangesInPasteboard), userInfo: nil, repeats: true)
    }

    @objc private func checkForChangesInPasteboard() {
        guard pasteboard.changeCount != changeCount else {
            return
        }

        exit(0)
    }
}


let pw = PasteboardWatcher()
pw.startListening()

// https://stackoverflow.com/questions/31944011/how-to-prevent-a-command-line-tool-from-exiting-before-asynchronous-operation-co
RunLoop.main.run()
