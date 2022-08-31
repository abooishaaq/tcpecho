import Control.Concurrent
import Control.Concurrent.Async
import Control.Monad
import Control.Monad.Loops
import Network.Simple.TCP
import Control.Exception.Safe

raceQuietly_ x y = race_ x y `catchIO` (\_ -> return ())

main = do
    bcast <- newChan
    forkIO $ forever $ readChan bcast
    serve (Host "127.0.0.1") "1337" $ \(s, addr) -> do
        bcast' <- dupChan bcast
        raceQuietly_
            (forever $ readChan bcast' >>= send s)
            (whileJust_ (recv s 4096) $ writeChan bcast')
