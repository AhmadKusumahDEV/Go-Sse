	defer func() {
			close(messageChan)
			messageChan = nil
		}()