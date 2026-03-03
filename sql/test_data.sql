insert into tags (id, name) values 
(1, 'muebles'),
(2, 'oficina'),
(3, 'cocina'),
(4, 'herramientas');

insert into articles (id, name, description) values
(1, 'escritorio 50x120 melamina', 'escritorio comprado en sodimac'),
(2, 'olla roja', 'se ve fea pero esta en buen estado. La herede y no la uso'),
(3, 'caladora black&decker alambrica', 'perfecto estado, con sus accesorios y caga. Se vende por upgrade'),
(4, 'Pantalla plana lcd viewsonic 4:3', 'viejita y pesada pero sirve para tenerla de pantalla auxiliar par un servidor, escritorio secundario, o algo por el estilo');

insert into articles_tags (article_id, tag_id) values
(1, 1),
(1, 2),
(2, 3),
(3, 4),
(4, 2);

insert into wishitems (id, name, external_url) values
(1, 'Cámara Web Pc Notebook Full Hd Usb Microfono 1080p', 'https://articulo.mercadolibre.cl/MLC-2973994090-camara-web-pc-notebook-full-hd-usb-microfono-1080p-_JM#polycard_client=recommendations_vip-pads-right&reco_backend=ranker_retrieval_system_ads&reco_model=rk_ent_v3_retsys_ads&reco_client=vip-pads-right&reco_item_pos=1&reco_backend_type=low_level&reco_id=37cce21a-fec0-4350-b3c1-781491d263d7&is_advertising=true&ad_domain=VIPCORE_RIGHT&ad_position=2&ad_click_id=M2M2ZDYyMDQtYzBmNC00ZTY2LTk3N2EtOTY5NDFmOTM5ZDAw'),
(2, 'Advantage360 Signature Series', 'https://kinesis-ergo.com/shop/advantage360-signature/');

insert into wishitems_tags (wishitem_id, tag_id) values
(1, 2),
(2, 1);
