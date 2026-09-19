(def init
  {:stage :first-load
   :mail []
   :loading #f})

(defn update [model action]
  (print action)
  (case (action :type)
        :first (-> model
                   (set :stage :alter)
                   (set :loading #t))
        :load-messages (if (action :status)
                         (-> model
                             (set :mail (action :value)))
                         model)
        model)
  )

(defn blank [message]
  (div {:class :blank}
   (span "No messages")))

(defn format-address [address]
  (let [name (take-while (partial /= "<") address)]
    [(take-while (partial /= "<") address)
     (->> address
          (drop-while (partial /= "<"))
          drop
          init
          string)
    ]))

(defn show-mail [m]
  (div {:class :mail}
       (let [address (format-address (m :address))]
         (div
          (span {:class :address-name}
                (head address))
          (span {:class :address-full}
                (-> address drop head))))
       (span {:class :subject}
             (m :subject))))

(defn show-inbox [model]
  (let [mail (model :mail)]
    (print mail (empty? mail))
    (if (empty? mail)
      (blank "No messages")
      (div {:class :inbox}
       ,(map show-mail mail)
       ,(map show-mail mail)
       ,(map show-mail mail)))))

(defn view [model]
  (print "view - model is" model)
  (div
   (h1 "Seamail")
   (show-inbox model)))

(defn on-view [model]
  (if (= (model :stage) :first-load)
  (do (send :first #f)
      (send-async (httpReq :load-messages
                          {:url "https://niliara.net/api/"})))
    ))
