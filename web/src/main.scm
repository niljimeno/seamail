(def init
  {:stage :first-load
   :mail []})

(defn update [model action]
  (print action)
  (case (action :type)
        :first (set model :stage :alter)
        :load-messages (do (print "fucking doing it")
                           (set model :mail (action :value)))
        model)
  )

(defn blank [message]
  (div {:class :blank}
   (span "No messages")))

(defn show-mail [m]
  (p (m :subject) (m :address)))

(defn show-inbox [model]
  (let [mail (model :mail)]
    (if (empty? mail)
      [(blank "No messages")]
      (map show-mail mail))))

(defn view [model]
  (print "draw" model)
  (div
   (h1 "Seamail")
   ,(show-inbox model)))

(defn on-view [model]
  (if (= (model :stage) :first-load)
  (do (send :first #f)
      (sendAsync (httpReq :load-messages
                          {:url "https://niliara.net/api/"})))
    ))
